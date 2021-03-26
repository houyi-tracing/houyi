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
	grpc2 "github.com/houyi-tracing/houyi/cmd/cs/app/handler/grpc"
	store2 "github.com/houyi-tracing/houyi/cmd/cs/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerParams struct {
	ListenPort int

	Logger *zap.Logger

	GossipRegistry gossip.Registry

	GossipSeed gossip.Seed

	ScaleFactor float64

	SST sst.SamplingStrategyTree

	TraceGraph tg.TraceGraph

	OperationStore store2.OperationStore

	Evaluator evaluator.Evaluator

	StrategyStore store2.StrategyStore

	MinSamplingRate float64
}

func StartGrpcServer(params *GrpcServerParams) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port:%d", params.ListenPort)
	}

	s := grpc.NewServer()
	if err = serverGrpc(s, lis, params); err != nil {
		return nil, err
	}

	return s, nil
}

func serverGrpc(s *grpc.Server, lis net.Listener, params *GrpcServerParams) error {
	rGrpcHandler := grpc2.NewRegistryGrpcHandler(params.Logger, params.GossipRegistry)
	api_v1.RegisterRegistryServer(s, rGrpcHandler)

	smGrpcHandler := grpc2.NewStrategyManagerGrpcHandler(
		params.Logger,
		params.SST,
		params.TraceGraph,
		params.OperationStore,
		params.Evaluator,
		params.ScaleFactor,
		params.MinSamplingRate,
		params.GossipSeed)
	api_v1.RegisterStrategyManagerServer(s, smGrpcHandler)

	params.Logger.Info("Starting gRPC server", zap.Int("port", params.ListenPort))
	go func() {
		if err := s.Serve(lis); err != nil {
			params.Logger.Fatal("Can not start gRPC server", zap.Error(err))
		}
	}()

	return nil
}
