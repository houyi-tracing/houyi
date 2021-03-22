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

package server

import (
	"fmt"
	"github.com/houyi-tracing/houyi/cmd/sm/app/handler"
	"github.com/houyi-tracing/houyi/cmd/sm/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerParams struct {
	Logger         *zap.Logger
	ListenPort     int
	ScaleFactor    *atomic.Float64
	StrategyStore  sst.SamplingStrategyTree
	TraceGraph     tg.TraceGraph
	OperationStore store.OperationStore
	Evaluator      evaluator.Evaluator
	GossipSeed     gossip.Seed
}

func StartGrpcServer(params *GrpcServerParams) (*grpc.Server, error) {
	params.Logger.Info("Starting gRPC server", zap.Int("port", params.ListenPort))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port: %d", params.ListenPort)
	}

	s := grpc.NewServer()
	if err := serveGrpc(s, lis, params); err != nil {
		params.Logger.Error("failed to server gRPC", zap.Error(err))
		return nil, err
	}

	return s, nil
}

func serveGrpc(s *grpc.Server, lis net.Listener, params *GrpcServerParams) error {
	grpcHandler := handler.NewGrpcHandler(
		params.Logger,
		params.StrategyStore,
		params.TraceGraph,
		params.OperationStore,
		params.Evaluator,
		params.GossipSeed,
		params.ScaleFactor,
	)

	api_v1.RegisterDynamicStrategyManagerServer(s, grpcHandler)
	api_v1.RegisterStrategyManagerServer(s, grpcHandler)
	api_v1.RegisterTraceGraphManagerServer(s, grpcHandler)
	api_v1.RegisterEvaluatorManagerServer(s, grpcHandler)

	go func() {
		if err := s.Serve(lis); err != nil {
			params.Logger.Fatal("failed to server grpc", zap.Error(err))
		}
	}()

	return nil
}
