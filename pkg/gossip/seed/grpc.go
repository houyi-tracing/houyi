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
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/yan-fuhai/go-ds/cache"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type grpcHandler struct {
	api_v1.UnimplementedSeedServer
	sync.Mutex

	logger   *zap.Logger
	msgCache cache.LRU
	seed     *seed
}

func newGrpcHandler(logger *zap.Logger, lruSize int, s *seed) *grpcHandler {
	return &grpcHandler{
		UnimplementedSeedServer: api_v1.UnimplementedSeedServer{},
		logger:                  logger,
		msgCache:                cache.NewLRU(lruSize),
		seed:                    s,
	}
}

// Sync send messages to randomly picked peers.
//
// For one message, every peer is in one of below three states:
//
// Susceptible(S): Initial state that a peer do not know this message. Peer in S would turn into R when it
// receive this message from another peer by probability k(probToR), otherwise turn to I.
//
// Infected(I): Middle state that a peer already knows this message and is mongering it. Peer in R would turn
// into R when it received this message again from other peers by probability k(probToR).
//
// Removed(R): Final state that a peer already knows the message and has quit spreading it. Peer in R
// would do nothing when it received this message again from other peers.
func (g *grpcHandler) Sync(_ context.Context, msg *api_v1.Message) (*api_v1.NullReply, error) {
	g.Lock()
	defer g.Unlock()

	s := g.seed

	s.logger.Debug("Received message", zap.String("msg", msg.String()))

	if cachedMsg := g.msgCache.Get(msg.MsgId); cachedMsg == nil {
		// Susceptible
		g.logger.Debug("Received new msg",
			zap.Int("node id", s.nodeId), zap.String("msg", msg.String()))

		switch msg.GetMsgType() {
		case api_v1.Message_NEW_RELATION:
			defer s.onNewRelation(msg.GetRelation())
		case api_v1.Message_EXPIRED_OPERATION:
			defer s.onExpiredOperation(msg.GetOperation())
		case api_v1.Message_EVALUATING_TAGS:
			defer s.onEvaluatingTags(msg.GetEvaluateTags())
		default:
			g.logger.Error("Unsupported type of message")
		}

		newItem := &msgCacheItem{
			state: Susceptible,
			msg:   msg,
		}

		rand.Seed(time.Now().UnixNano())
		if rand.Float64() < s.probToR {
			g.logger.Debug("Switch to R from S", zap.Int64("msg ID", msg.MsgId))
			newItem.state = Removed
		} else {
			g.logger.Debug("Switch to I from S", zap.Int64("msg ID", msg.MsgId))
			newItem.state = Infected
		}
		g.msgCache.Put(msg.MsgId, newItem)

		s.msgSender <- msg
		g.logger.Debug("Sent message", zap.Int64("msg ID", msg.MsgId), zap.String("state", "S"))
	} else if mci, ok := cachedMsg.(*msgCacheItem); ok {
		if mci.state == Infected {
			// Infected
			rand.Seed(time.Now().UnixNano())
			if rand.Float64() < s.probToR {
				g.logger.Debug("Switch to R from I", zap.Int64("msg ID", msg.MsgId))
				mci.state = Removed
			} else {
				s.msgSender <- msg
				g.logger.Debug("Sent message", zap.Int64("msg ID", msg.MsgId), zap.String("state", "I"))
			}
		} else {
			// Removed
			// do nothing
		}
	}

	return &api_v1.NullReply{}, nil
}
