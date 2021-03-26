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

package registry

import (
	"github.com/bwmarrin/snowflake"
	grpc2 "github.com/houyi-tracing/houyi/cmd/cs/app/handler/grpc"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	maxNode int64 = -1 ^ (-1 << snowflake.NodeBits)
)

type registry struct {
	logger *zap.Logger

	randomPick        int
	probToR           float64
	heartbeatInterval time.Duration

	peers       SeedSet
	grpcHandler *grpc2.RegistryGrpcHandler
	stop        chan *sync.WaitGroup
}

func NewRegistry(logger *zap.Logger, randomPick int, probToR float64, hbIntvl time.Duration) gossip.Registry {
	r := &registry{
		logger:            logger,
		randomPick:        randomPick,
		probToR:           probToR,
		heartbeatInterval: hbIntvl,
		peers:             NewSeedSet(),
		stop:              make(chan *sync.WaitGroup),
	}
	return r
}

func (r *registry) Start() error {
	r.logger.Info("Starting registry",
		zap.Int("random pick", r.randomPick),
		zap.Float64("probability to R", r.probToR),
		zap.Duration("heartbeat interval", r.heartbeatInterval))

	go r.timer(r.removeDeadNodes, r.heartbeatInterval)
	return nil
}

func (r *registry) Stop() error {
	var wg sync.WaitGroup
	wg.Add(1)
	r.stop <- &wg
	wg.Done()
	return nil
}

func (r *registry) AllSeeds() []*api_v1.Peer {
	return r.peers.AllSeeds()
}

func (r *registry) Register(ip string, port int) (int64, int, time.Duration, float64) {
	newNodeId := r.peers.Add(ip, port)
	r.logger.Info("Received register request from seed",
		zap.String("IP", ip),
		zap.Int("port", port),
		zap.Int64("node id", newNodeId))
	return newNodeId, r.randomPick, r.heartbeatInterval, r.probToR
}

func (r *registry) Heartbeat(id int64, ip string, port int) (int64, []*api_v1.Peer) {
	node := r.peers.GetNode(id)
	if !r.peers.Has(id) || node == nil || node.Ip != ip || node.Port != int64(port) {
		// The seed id of registered seed would be recycled because it has not sent a heartbeat for a long time,
		// or the configuration (IP, port etc) of seed has changed.
		// In those cases, we should generate new unique ID for this seed because old ID may has been sign to another
		// new seed.
		id = r.peers.Add(ip, port)
	} else {
		r.peers.Refresh(id)
	}
	allPeers := r.peers.AllPeers(id) // exclude the node sent this request
	return id, allPeers
}

func (r *registry) timer(f func(), interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f()
		case wg := <-r.stop:
			wg.Done()
			return
		}
	}
}

func (r *registry) removeDeadNodes() {
	for _, id := range r.peers.AllIds() {
		if r.peers.IsDead(id, r.heartbeatInterval) {
			r.logger.Info("Removed dead seed",
				zap.String("IP", r.peers.GetNode(id).Ip),
				zap.Int64("port", r.peers.GetNode(id).Port),
				zap.Int64("node id", id))
			r.peers.Remove(id)
		}
	}
}
