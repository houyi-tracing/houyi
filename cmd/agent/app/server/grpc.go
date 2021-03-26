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
	"github.com/houyi-tracing/houyi/cmd/agent/app/handler"
	"github.com/houyi-tracing/houyi/cmd/agent/app/transport"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/routing"
	jaeger "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerParams struct {
	Logger               *zap.Logger
	ListenPort           int
	CollectorEndpoint    *routing.Endpoint
	ConfigServerEndpoint *routing.Endpoint
}

func StartGrpcServer(params *GrpcServerParams) (*grpc.Server, error) {
	params.Logger.Info("Starting gRPC server", zap.Int("port", params.ListenPort))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	if err := serveGrpc(s, lis, params); err != nil {
		return nil, err
	} else {
		return s, nil
	}
}

func serveGrpc(s *grpc.Server, lis net.Listener, params *GrpcServerParams) error {
	cTransport := transport.NewCollectorTransport(params.Logger, params.CollectorEndpoint)
	smTransport := transport.NewStrategyManagerTransport(params.Logger, params.ConfigServerEndpoint)

	h := handler.NewGrpcHandler(params.Logger, cTransport, smTransport)

	api_v1.RegisterStrategyManagerServer(s, h)
	jaeger.RegisterCollectorServiceServer(s, h)

	go func() {
		if err := s.Serve(lis); err != nil {
			params.Logger.Fatal("Failed to server gRPC", zap.Error(err))
		}
	}()

	return nil
}
