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
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerParams struct {
	GrpcHandler    api_v1.RegistryServer
	ListenPort     int
	Logger         *zap.Logger
	GossipRegistry gossip.Registry
}

func StartGrpcServer(params *GrpcServerParams) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", params.ListenPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port:%d", params.ListenPort)
	}

	s := grpc.NewServer()
	if err := serverGrpc(s, lis, params); err != nil {
		return nil, err
	}

	return s, nil
}

func serverGrpc(s *grpc.Server, lis net.Listener, params *GrpcServerParams) error {
	api_v1.RegisterRegistryServer(s, params.GrpcHandler)

	params.Logger.Info("Starting gRPC server", zap.Int("port", params.ListenPort))
	go func() {
		if err := s.Serve(lis); err != nil {
			params.Logger.Fatal("Can not launch gRPC server", zap.Error(err))
		}
	}()

	return nil
}
