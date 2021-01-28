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
	"github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/cmd/collector/app/processor"
	"github.com/houyi-tracing/houyi/cmd/collector/app/server"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CollectorParams struct {
	Logger         *zap.Logger
	GrpcListenPort int
	TraceGraph     tg.TraceGraph
	GossipSeed     gossip.Seed
	Evaluator      evaluator.Evaluator
	SpanProcessor  processor.SpanProcessor
	SpanFilter     filter.SpanFilter
	SpanWriter     spanstore.Writer
}

type Collector struct {
	logger         *zap.Logger
	grpcServer     *grpc.Server
	grpcListenPort int

	spanProcessor processor.SpanProcessor
	spanFilter    filter.SpanFilter
	spanWriter    spanstore.Writer

	traceGraph tg.TraceGraph
	evaluator  evaluator.Evaluator
	gossipSeed gossip.Seed
}

func NewCollector(params *CollectorParams) *Collector {
	return &Collector{
		logger:         params.Logger,
		spanProcessor:  params.SpanProcessor,
		spanFilter:     params.SpanFilter,
		spanWriter:     params.SpanWriter,
		traceGraph:     params.TraceGraph,
		evaluator:      params.Evaluator,
		gossipSeed:     params.GossipSeed,
		grpcListenPort: params.GrpcListenPort,
	}
}

func (c *Collector) Start() error {
	if gS, err := server.StartGrpcServer(&server.GrpcServerParams{
		Logger:        c.logger,
		ListenPort:    c.grpcListenPort,
		SpanProcessor: c.spanProcessor,
	}); err != nil {
		return err
	} else {
		c.grpcServer = gS
	}

	return nil
}

func (c *Collector) Close() error {
	c.grpcServer.GracefulStop()
	return nil
}
