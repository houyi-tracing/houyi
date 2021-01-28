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

package app

import (
	"github.com/houyi-tracing/houyi/cmd/agent/app/server"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AgentParams struct {
	Logger                  *zap.Logger
	GrpcListenPort          int
	CollectorEndpoint       *routing.Endpoint
	StrategyManagerEndpoint *routing.Endpoint
}

// Agent is used to mask the routing information of the collector and strategy manager for the client.
type Agent struct {
	logger *zap.Logger

	grpcListenPort int
	grpcServer     *grpc.Server

	cEp  *routing.Endpoint // endpoint of collector
	smEp *routing.Endpoint // endpoint of strategy manager
}

func NewAgent(params *AgentParams) *Agent {
	return &Agent{
		logger:         params.Logger,
		cEp:            params.CollectorEndpoint,
		smEp:           params.StrategyManagerEndpoint,
		grpcListenPort: params.GrpcListenPort,
	}
}

func (a *Agent) Start() error {
	if grpcServer, err := server.StartGrpcServer(&server.GrpcServerParams{
		Logger:                  a.logger,
		ListenPort:              a.grpcListenPort,
		CollectorEndpoint:       a.cEp,
		StrategyManagerEndpoint: a.smEp,
	}); err != nil {
		return err
	} else {
		a.grpcServer = grpcServer
	}

	return nil
}

func (a *Agent) Stop() error {
	a.grpcServer.GracefulStop()
	return nil
}
