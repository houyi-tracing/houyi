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

package handler

import (
	"context"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	api_v1.UnimplementedRegistryServer

	logger   *zap.Logger
	registry gossip.Registry
}

func NewGrpcHandler(logger *zap.Logger, registry gossip.Registry) api_v1.RegistryServer {
	return &GrpcHandler{
		logger:   logger,
		registry: registry,
	}
}

func (h *GrpcHandler) Register(ctx context.Context, req *api_v1.RegisterRequest) (*api_v1.RegisterRely, error) {
	ip := req.GetIp()
	port := req.GetPort()

	nodeId, randomPick, interval, probToR := h.registry.Register(ip, int(port))
	return &api_v1.RegisterRely{
		NodeId:     nodeId,
		Interval:   interval.Nanoseconds() * 2 / 3,
		RandomPick: int64(randomPick),
		ProbToR:    probToR,
	}, nil
}

func (h *GrpcHandler) Heartbeat(ctx context.Context, req *api_v1.HeartbeatRequest) (*api_v1.HeartbeatReply, error) {
	id := req.GetNodeId()
	ip := req.GetIp()
	port := req.GetPort()

	id, peers := h.registry.Heartbeat(id, ip, int(port))
	return &api_v1.HeartbeatReply{
		NodeId: id,
		Peers:  peers,
	}, nil
}
