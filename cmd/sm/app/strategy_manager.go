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
	"github.com/houyi-tracing/houyi/cmd/sm/app/server"
	"github.com/houyi-tracing/houyi/cmd/sm/app/store"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type StrategyManagerParams struct {
	Logger          *zap.Logger
	GrpcListenPort  int
	ScaleFactor     *atomic.Float64
	RefreshInterval time.Duration
	StrategyStore   sst.SamplingStrategyTree
	TraceGraph      tg.TraceGraph
	Evaluator       evaluator.Evaluator
	GossipSeed      gossip.Seed
}

type StrategyManager struct {
	sync.RWMutex

	scaleFactor *atomic.Float64

	grpcListenPort  int
	refreshInterval time.Duration

	logger     *zap.Logger
	grpcServer *grpc.Server

	strategyStore sst.SamplingStrategyTree
	tg            tg.TraceGraph
	eval          evaluator.Evaluator
	seed          gossip.Seed
	opStore       store.OperationStore
}

func NewStrategyManager(params *StrategyManagerParams) *StrategyManager {
	ret := &StrategyManager{
		logger:          params.Logger,
		grpcListenPort:  params.GrpcListenPort,
		refreshInterval: params.RefreshInterval,
		strategyStore:   params.StrategyStore,
		tg:              params.TraceGraph,
		eval:            params.Evaluator,
		seed:            params.GossipSeed,
		scaleFactor:     params.ScaleFactor,
	}
	ret.opStore = store.NewOperationStore(params.Logger, ret.refreshInterval, ret.seed)
	return ret
}

func (m *StrategyManager) Start() error {
	if grpcServer, err := server.StartGrpcServer(&server.GrpcServerParams{
		Logger:         m.logger,
		ListenPort:     m.grpcListenPort,
		StrategyStore:  m.strategyStore,
		TraceGraph:     m.tg,
		OperationStore: m.opStore,
		Evaluator:      m.eval,
		GossipSeed:     m.seed,
		ScaleFactor:    m.scaleFactor,
	}); err != nil {
		return err
	} else {
		m.grpcServer = grpcServer
	}

	m.opStore.Start()

	return nil
}

func (m *StrategyManager) Stop() error {
	m.grpcServer.GracefulStop()
	m.opStore.Stop()
	if err := m.seed.Stop(); err != nil {
		return err
	}
	return nil
}
