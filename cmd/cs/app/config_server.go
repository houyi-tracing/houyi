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
	"github.com/houyi-tracing/houyi/cmd/cs/app/server"
	"github.com/houyi-tracing/houyi/cmd/cs/app/store"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ConfigurationServerParams struct {
	Logger *zap.Logger

	GrpcListenPort int

	HttpListenPort int

	GossipSeed gossip.Seed

	GossipRegistry gossip.Registry

	TraceGraph tg.TraceGraph

	Evaluator evaluator.Evaluator

	StrategyStore store.StrategyStore

	SST sst.SamplingStrategyTree

	OperationStore store.OperationStore

	ScaleFactor float64

	MinSamplingRate float64
}

type ConfigurationServer struct {
	logger *zap.Logger

	grpcServer *grpc.Server

	grpcListenPort int

	httpListenPort int

	gossipSeed gossip.Seed

	gossipRegistry gossip.Registry

	traceGraph tg.TraceGraph

	evaluator evaluator.Evaluator

	sst sst.SamplingStrategyTree

	strategyStore store.StrategyStore

	opStore store.OperationStore

	scaleFactor float64

	minSamplingRate float64
}

func NewConfigServer(params *ConfigurationServerParams) *ConfigurationServer {
	return &ConfigurationServer{
		logger:          params.Logger,
		gossipSeed:      params.GossipSeed,
		gossipRegistry:  params.GossipRegistry,
		traceGraph:      params.TraceGraph,
		evaluator:       params.Evaluator,
		strategyStore:   params.StrategyStore,
		opStore:         params.OperationStore,
		grpcListenPort:  params.GrpcListenPort,
		httpListenPort:  params.HttpListenPort,
		scaleFactor:     params.ScaleFactor,
		sst:             params.SST,
		minSamplingRate: params.MinSamplingRate,
	}
}

func (cs *ConfigurationServer) Start() error {
	var err error

	if cs.grpcServer, err = server.StartGrpcServer(&server.GrpcServerParams{
		ListenPort:      cs.grpcListenPort,
		Logger:          cs.logger,
		GossipRegistry:  cs.gossipRegistry,
		GossipSeed:      cs.gossipSeed,
		ScaleFactor:     cs.scaleFactor,
		SST:             cs.sst,
		TraceGraph:      cs.traceGraph,
		OperationStore:  cs.opStore,
		Evaluator:       cs.evaluator,
		StrategyStore:   cs.strategyStore,
		MinSamplingRate: cs.minSamplingRate,
	}); err != nil {
		return err
	}

	if err = server.StartHttpServer(&server.HttpServerParams{
		ListenPort:     cs.httpListenPort,
		Logger:         cs.logger,
		TraceGraph:     cs.traceGraph,
		StrategyStore:  cs.strategyStore,
		Evaluator:      cs.evaluator,
		GossipRegistry: cs.gossipRegistry,
	}); err != nil {
		return err
	}

	if err = cs.gossipSeed.Start(); err != nil {
		return err
	}

	cs.opStore.Start()

	return nil
}

func (cs *ConfigurationServer) Stop() error {
	_ = cs.gossipSeed.Stop()
	cs.opStore.Stop()
	cs.grpcServer.GracefulStop()
	return nil
}
