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
	"github.com/houyi-tracing/houyi/cmd/collector/app/handler"
	"github.com/houyi-tracing/houyi/cmd/collector/app/processor"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerParams struct {
	Logger        *zap.Logger
	ListenPort    int
	SpanProcessor processor.SpanProcessor
	Evaluator     evaluator.Evaluator
}

func StartGrpcServer(params *GrpcServerParams) (*grpc.Server, error) {
	params.Logger.Info("Starting gRPC server", zap.Int("port", params.ListenPort))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	if err := serveGrpc(server, lis, params); err != nil {
		return nil, err
	} else {
		return server, nil
	}
}

func serveGrpc(server *grpc.Server, lis net.Listener, params *GrpcServerParams) error {
	gh := handler.NewGrpcHandler(params.Logger, params.Evaluator, params.SpanProcessor)

	api_v2.RegisterCollectorServiceServer(server, gh)
	api_v1.RegisterEvaluatorManagerServer(server, gh)

	go func() {
		if err := server.Serve(lis); err != nil {
			params.Logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	return nil
}
