// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
// Copyright (c) 2020 The Houyi Authors.
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
	handler2 "github.com/houyi-tracing/houyi/cmd/collector/app/handler"
	"github.com/houyi-tracing/houyi/cmd/collector/app/sampling"
	"github.com/jaegertracing/jaeger/pkg/healthcheck"
	"os"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/collector/app/handler"
	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	zs "github.com/jaegertracing/jaeger/cmd/collector/app/sanitizer/zipkin"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// SpanHandlerBuilder holds configuration required for handlers
type SpanHandlerBuilder struct {
	SpanWriter     spanstore.Writer
	CollectorOpts  CollectorOptions
	Logger         *zap.Logger
	MetricsFactory metrics.Factory
	StrategyStore  sampling.AdaptiveStrategyStore
	SpanFilter     filter.SpanFilter
	HealthCheck    *healthcheck.HealthCheck
}

// SpanHandlers holds instances to the span handlers built by the SpanHandlerBuilder
type SpanHandlers struct {
	ZipkinSpansHandler   handler.ZipkinSpansHandler
	JaegerBatchesHandler handler.JaegerBatchesHandler
	GRPCHandler          *handler.GRPCHandler
	HttpHandler          *handler2.APIHandler
}

// BuildSpanProcessor builds the span processor to be used with the handlers
func (b *SpanHandlerBuilder) BuildSpanProcessor() processor.SpanProcessor {
	hostname, _ := os.Hostname()
	svcMetrics := b.metricsFactory()
	hostMetrics := svcMetrics.Namespace(metrics.NSOptions{Tags: map[string]string{"host": hostname}})

	return NewAdaptiveSamplingSpanProcessor(
		b.SpanWriter,
		Options.ServiceMetrics(svcMetrics),
		Options.HostMetrics(hostMetrics),
		Options.Logger(b.logger()),
		Options.SpanFilter(defaultSpanFilter),
		Options.NumWorkers(b.CollectorOpts.NumWorkers),
		Options.QueueSize(b.CollectorOpts.QueueSize),
		Options.CollectorTags(b.CollectorOpts.CollectorTags),
		Options.DynQueueSizeWarmup(uint(b.CollectorOpts.QueueSize)), // same as queue size for now
		Options.DynQueueSizeMemory(b.CollectorOpts.DynQueueSizeMemory),
		// adaptive sampling options place from here
		Options.LruCapacity(b.CollectorOpts.LruCapacity),
		Options.MaxRetries(b.CollectorOpts.MaxRetries),
		Options.StoreRefreshInterval(b.CollectorOpts.StoreRefreshInterval),
		Options.AdaptiveSamplingSpanFilter(b.SpanFilter),
		Options.StrategyStore(b.StrategyStore),
		Options.RetryQueueNumWorkers(b.CollectorOpts.RetryQueueNumWorkers),
	)
}

// BuildHandlers builds span handlers (Zipkin, Jaeger)
func (b *SpanHandlerBuilder) BuildHandlers(spanProcessor processor.SpanProcessor) *SpanHandlers {
	jaegerBatchesHandler := handler.NewJaegerSpanHandler(b.Logger, spanProcessor)
	return &SpanHandlers{
		handler.NewZipkinSpanHandler(b.Logger, spanProcessor, zs.NewChainedSanitizer(zs.StandardSanitizers...)),
		jaegerBatchesHandler,
		handler.NewGRPCHandler(b.Logger, spanProcessor),
		handler2.NewAPIHandler(jaegerBatchesHandler, b.StrategyStore, b.SpanFilter, b.HealthCheck),
	}
}

func defaultSpanFilter(*model.Span) bool {
	return true
}

func (b *SpanHandlerBuilder) logger() *zap.Logger {
	if b.Logger == nil {
		return zap.NewNop()
	}
	return b.Logger
}

func (b *SpanHandlerBuilder) metricsFactory() metrics.Factory {
	if b.MetricsFactory == nil {
		return metrics.NullFactory
	}
	return b.MetricsFactory
}
