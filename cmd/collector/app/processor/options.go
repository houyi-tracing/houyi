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

package processor

import (
	"github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

type options struct {
	numWorkers       int
	registryEndpoint *routing.Endpoint

	filterSpan              filter.FilterSpan
	evaluateSpan            evaluator.EvaluateSpan
	spanWriter              spanstore.Writer
	traceGraph              tg.TraceGraph
	seed                    gossip.Seed
	strategyManagerEndpoint *routing.Endpoint
}

var Options options

type Option func(opt *options)

func (options) NumWorkers(n int) Option {
	return func(opt *options) {
		opt.numWorkers = n
	}
}

func (options) FilterSpan(f filter.FilterSpan) Option {
	return func(opt *options) {
		opt.filterSpan = f
	}
}

func (options) EvaluateSpan(f evaluator.EvaluateSpan) Option {
	return func(opt *options) {
		opt.evaluateSpan = f
	}
}

func (options) SpanWriter(w spanstore.Writer) Option {
	return func(opt *options) {
		opt.spanWriter = w
	}
}

func (options) StrategyManagerEndpoint(ep *routing.Endpoint) Option {
	return func(opt *options) {
		opt.strategyManagerEndpoint = ep
	}
}

func (options) GossipSeed(seed gossip.Seed) Option {
	return func(opt *options) {
		opt.seed = seed
	}
}

func (options) TraceGraph(tg tg.TraceGraph) Option {
	return func(opt *options) {
		opt.traceGraph = tg
	}
}

func (o *options) apply(opts ...Option) *options {
	for _, op := range opts {
		op(o)
	}

	if o.numWorkers == 0 {
		o.numWorkers = DefaultNumWorkers
	}
	return o
}
