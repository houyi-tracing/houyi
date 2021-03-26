// Copyright (c) 2021 The Houyi Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package seed

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	Susceptible = iota
	Infected
	Removed
)

type msgCacheItem struct {
	state int
	msg   *api_v1.Message
}

type seed struct {
	options

	lock              sync.Mutex
	nodeId            int
	randomPick        int
	probToR           float64
	heartbeatInterval time.Duration
	msgIdGenerator    *snowflake.Node
	grpcServer        *grpc.Server
	grpcHandler       api_v1.SeedServer
	logger            *zap.Logger
	peers             []*api_v1.Peer
	msgSender         chan *api_v1.Message
	stopMsgSender     chan *sync.WaitGroup
	stopTimer         chan *sync.WaitGroup
}

func NewSeed(logger *zap.Logger, opts ...Option) gossip.Seed {
	s := &seed{
		lock:          sync.Mutex{},
		options:       Options.apply(opts...),
		logger:        logger,
		peers:         make([]*api_v1.Peer, 0),
		stopMsgSender: make(chan *sync.WaitGroup),
		stopTimer:     make(chan *sync.WaitGroup),
	}
	s.msgSender = make(chan *api_v1.Message, s.lruSize/4)
	return s
}

func (s *seed) OnNewRelation(f func(rel *api_v1.Relation)) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.onNewRelation = f
}

func (s *seed) OnExpiredOperation(f func(op *api_v1.Operation)) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.onExpiredOperation = f
}

func (s *seed) OnEvaluatingTags(f func(tags *api_v1.EvaluatingTags)) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.onEvaluatingTags = f
}

func (s *seed) OnNewOperation(f func(op *api_v1.Operation)) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.onNewOperation = f
}

func (s *seed) MongerEvaluatingTags(tags *api_v1.EvaluatingTags) {
	msg := &api_v1.Message{
		MsgId:   s.msgIdGenerator.Generate().Int64(),
		MsgType: api_v1.Message_EVALUATING_TAGS,
		Msg: &api_v1.Message_EvaluateTags{
			EvaluateTags: tags,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if s.grpcHandler != nil {
		_, _ = s.grpcHandler.Sync(ctx, msg)
	} else {
		s.logger.Error("Grpc handler does not ready")
	}
}

func (s *seed) MongerExpiredOperation(op *api_v1.Operation) {
	msg := &api_v1.Message{
		MsgId:   s.msgIdGenerator.Generate().Int64(),
		MsgType: api_v1.Message_EXPIRED_OPERATION,
		Msg: &api_v1.Message_Operation{
			Operation: op,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if s.grpcHandler != nil {
		_, _ = s.grpcHandler.Sync(ctx, msg)
	} else {
		s.logger.Error("Grpc handler does not ready")
	}
}

func (s *seed) MongerNewRelation(rel *api_v1.Relation) {
	msg := &api_v1.Message{
		MsgId:   s.msgIdGenerator.Generate().Int64(),
		MsgType: api_v1.Message_NEW_RELATION,
		Msg: &api_v1.Message_Relation{
			Relation: rel,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if s.grpcHandler != nil {
		_, _ = s.grpcHandler.Sync(ctx, msg)
	} else {
		s.logger.Error("Grpc handler does not ready")
	}
}

func (s *seed) MongerNewOperation(op *api_v1.Operation) {
	msg := &api_v1.Message{
		MsgId:   s.msgIdGenerator.Generate().Int64(),
		MsgType: api_v1.Message_NEW_OPERATION,
		Msg: &api_v1.Message_Operation{
			Operation: op,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if s.grpcHandler != nil {
		_, _ = s.grpcHandler.Sync(ctx, msg)
	} else {
		s.logger.Error("Grpc handler does not ready")
	}
}

func (s *seed) Start() error {
	// get global item id from registry
	if err := s.register(); err != nil {
		return err
	}

	// send heartbeat message to registry to get peers from registry
	if err := s.heartbeat(); err != nil {
		return err
	}

	if idGenerator, err := snowflake.NewNode(int64(s.nodeId)); err != nil {
		return err
	} else {
		s.msgIdGenerator = idGenerator
	}

	go s.startGrpcServer()
	go s.msgMonger()
	go s.timer(s.heartbeat, s.heartbeatInterval)
	return nil
}

func (s *seed) timer(f func() error, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := f(); err != nil {
				s.logger.Error("Failed to run function in ticker", zap.Error(err))
			}
		case wg := <-s.stopTimer:
			wg.Done()
			return
		}
	}
}

func (s *seed) Stop() error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	var wg sync.WaitGroup
	wg.Add(2)
	s.stopTimer <- &wg
	s.stopMsgSender <- &wg
	wg.Wait()
	return nil
}

func (s *seed) startGrpcServer() {
	s.logger.Info("Starting gGRPC server for gossip seed",
		zap.Int("port", s.listenPort),
		zap.String("configuration server endpoint", s.configServerEp.String()),
		zap.Int("LRU size for message cache", s.lruSize))

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.listenPort))
	if err != nil {
		s.logger.Fatal("Failed to listen tcp for seed", zap.Error(err))
	}

	s.grpcServer = grpc.NewServer()
	s.grpcHandler = newGrpcHandler(s.logger, s.lruSize, s)
	api_v1.RegisterSeedServer(s.grpcServer, s.grpcHandler)
	if err := s.grpcServer.Serve(conn); err != nil {
		s.logger.Fatal("Failed to server grpc for seed", zap.Error(err))
	}
}

// sendMsg send message to one item.
func (s *seed) sendMsg(ip string, port int64, msg *api_v1.Message) {
	conn, err := grpc.Dial(formatEndpoint(ip, port), grpc.WithInsecure(), grpc.WithBlock())
	if conn == nil || err != nil {
		s.logger.Fatal("Could not dial remote seed", zap.String("ip", ip), zap.Int64("port", port))
	} else {
		defer conn.Close()
	}

	c := api_v1.NewSeedClient(conn)
	_, err = c.Sync(context.TODO(), msg)
	if err != nil {
		s.logger.Error("Failed to monger message", zap.String("IP", ip), zap.Int64("port", port))
	}
}

// msgMonger sends message to randomly picked peers.
func (s *seed) msgMonger() {
	for {
		select {
		case msg := <-s.msgSender:
			s.lock.Lock()
			picks := randomlyPick(len(s.peers), s.randomPick)
			for _, peerIdx := range picks {
				s.sendMsg(s.peers[peerIdx].GetIp(), s.peers[peerIdx].GetPort(), msg)
			}
			s.lock.Unlock()
		case wg := <-s.stopMsgSender:
			wg.Done()
			return
		}
	}
}

func (s *seed) register() error {
	var conn *grpc.ClientConn
	var err error
	for {
		conn, err = grpc.Dial(s.configServerEp.String(), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil || conn == nil {
			s.logger.Error("failed to dial to registry")
			time.Sleep(time.Second * 5)
			continue
		}
		defer conn.Close()

		c := api_v1.NewRegistryClient(conn)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ip, err := getClientIp()
		if err != nil {
			return err
		}
		req := &api_v1.RegisterRequest{
			Ip:   ip,
			Port: int64(s.listenPort),
		}

		reply := &api_v1.RegisterRely{}
		if reply, err = c.Register(ctx, req); err != nil {
			return nil
		} else {
			s.nodeId = int(reply.NodeId)
			s.randomPick = int(reply.RandomPick)
			s.heartbeatInterval = time.Duration(reply.Interval)
			s.probToR = reply.ProbToR

			s.logger.Info("Received reply from registry",
				zap.Int("node id", s.nodeId),
				zap.Int("random pick", s.randomPick),
				zap.Float64("probability to R", s.probToR),
				zap.Duration("heartbeat interval", s.heartbeatInterval))
			return nil
		}
	}
}

func (s *seed) heartbeat() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var conn *grpc.ClientConn
	var err error
	if conn, err = grpc.Dial(s.configServerEp.String(), grpc.WithInsecure()); conn == nil || err != nil {
		return err
	}
	defer conn.Close()

	c := api_v1.NewRegistryClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ip, err := getClientIp()
	if err != nil {
		return err
	}
	req := &api_v1.HeartbeatRequest{
		Ip:     ip,
		NodeId: int64(s.nodeId),
		Port:   int64(s.listenPort),
	}

	reply := &api_v1.HeartbeatReply{}
	if reply, err = c.Heartbeat(ctx, req); err != nil {
		return err
	} else {
		//s.logger.Debug("Received heartbeat reply from registry",
		//	zap.Int64("node id", reply.NodeId),
		//	zap.Int("peers", len(reply.Peers)))
		s.peers = reply.Peers
		if s.nodeId != int(reply.NodeId) {
			s.logger.Debug("Received new node id from registry", zap.Int64("node id", reply.NodeId))
		}
		s.nodeId = int(reply.NodeId)
		return nil
	}
}

// randomlyPick randomly picks n gossip peers
func randomlyPick(total, need int) []int {
	need = min(total, need)
	ret := make([]int, 0, need)
	rand.Seed(time.Now().UnixNano())
	if total != 0 {
		picked := make(map[int]bool)
		for len(picked) < need {
			pick := rand.Intn(total)
			if _, isPicked := picked[pick]; !isPicked {
				picked[pick] = true
			}
		}
		for k := range picked {
			ret = append(ret, k)
		}
	}
	return ret
}

func formatEndpoint(ip string, port int64) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("Can not find the client ip address!")

}
