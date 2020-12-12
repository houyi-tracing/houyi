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
	"github.com/houyi-tracing/houyi/cmd/collector/app/sampling"
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/jaegertracing/jaeger/cmd/collector/app"
	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/queue"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/yan-fuhai/go-ds/cache"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"time"
)

// AdaptiveSamplingSpanProcessor overrides spanProcessor for adaptive sampling.
type AdaptiveSamplingSpanProcessor struct {
	*spanProcessor

	maxRetries           int
	retryQueueNumWorkers int
	store                sampling.AdaptiveStrategyStore
	lru                  cache.LRU // key: spanID, value: ExecutionGraphNode
	retryQueue           *queue.BoundedQueue
	retryQueueItemExpire time.Duration

	// ATTENTION: spanFilter is different from the filterSpan of spanProcessor.
	// This one filters spans that meets conditions for increasing the sampling probability of them.
	// The other one is for rejecting spans that failed to meet the conditions.
	spanFilter filter.SpanFilter
}

// The retry queue stores those spans that failed to update trace graph due to delay of parent spans.
// A span would NOT be pushed to retry queue when it's expired.
type retryQueueItem struct {
	span     *model.Span
	parentID model.SpanID
	since    time.Time
}

func NewAdaptiveSamplingSpanProcessor(spanWriter spanstore.Writer, opts ...Option) processor.SpanProcessor {
	sp := newAdaptiveSamplingSpanProcessor(spanWriter, opts...)

	sp.queue.StartConsumers(sp.numWorkers, func(item interface{}) {
		if qItem, ok := item.(*queueItem); ok {
			sp.processItemFromQueue(qItem)
		} else {
			sp.logger.Error("failed to convert to queueItem")
		}
	})

	sp.retryQueue.StartConsumers(sp.retryQueueNumWorkers, func(item interface{}) {
		if rQItem, ok := item.(*retryQueueItem); ok {
			sp.processItemFromRetryQueue(rQItem)
		} else {
			sp.logger.Error("failed to convert to retryQueueItem")
		}
	})

	sp.background(1*time.Second, sp.updateGauges)

	if sp.dynQueueSizeMemory > 0 {
		sp.background(1*time.Minute, sp.updateQueueSize)
	}

	return sp
}

func newAdaptiveSamplingSpanProcessor(spanWriter spanstore.Writer, opts ...Option) *AdaptiveSamplingSpanProcessor {
	options := Options.apply(opts...)
	handlerMetrics := app.NewSpanProcessorMetrics(
		options.serviceMetrics,
		options.hostMetrics,
		options.extraFormatTypes)
	droppedItemHandler := func(item interface{}) {
		handlerMetrics.SpansDropped.Inc(1)
	}
	boundedQueue := queue.NewBoundedQueue(options.queueSize, droppedItemHandler)

	sp := spanProcessor{
		queue:              boundedQueue,
		metrics:            handlerMetrics,
		logger:             options.logger,
		preProcessSpans:    options.preProcessSpans,
		filterSpan:         options.spanFilter,
		sanitizer:          options.sanitizer,
		reportBusy:         options.reportBusy,
		numWorkers:         options.numWorkers,
		spanWriter:         spanWriter,
		collectorTags:      options.collectorTags,
		stopCh:             make(chan struct{}),
		dynQueueSizeMemory: options.dynQueueSizeMemory,
		dynQueueSizeWarmup: options.dynQueueSizeWarmup,
		bytesProcessed:     atomic.NewUint64(0),
		spansProcessed:     atomic.NewUint64(0),
	}

	retryQueueDroppedItemHandler := func(item interface{}) {
		retryQItem := item.(*retryQueueItem)
		sp.logger.Info("retry queue dropped span", zap.String("operation name", retryQItem.span.OperationName))
	}
	retryQueue := queue.NewBoundedQueue(options.queueSize, retryQueueDroppedItemHandler)
	adsp := AdaptiveSamplingSpanProcessor{
		retryQueue:           retryQueue,
		retryQueueNumWorkers: options.retryQueueNumWorkers,
		retryQueueItemExpire: options.retryQueueItemExpire,
		spanFilter:           options.asSpanFilter,
		spanProcessor:        &sp,
		lru:                  cache.NewLRU(options.lruCapacity),
		store:                options.store,
	}

	processSpanFuncs := []app.ProcessSpan{options.preSave, adsp.updateTraceGraph, adsp.filterSpanToPromote, adsp.saveSpan}
	if options.dynQueueSizeMemory > 0 {
		// add to processSpanFuncs
		options.logger.Info("Dynamically adjusting the queue size at runtime.",
			zap.Uint("memory-mib", options.dynQueueSizeMemory/1024/1024),
			zap.Uint("queue-size-warmup", options.dynQueueSizeWarmup))
		processSpanFuncs = append(processSpanFuncs, adsp.countSpan)
	}

	sp.background(options.storeRefreshInterval, adsp.store.RemoveExpired)
	adsp.processSpan = app.ChainedProcessSpan(processSpanFuncs...)
	return &adsp
}

func (sp *AdaptiveSamplingSpanProcessor) UpdateFilter(newFilter filter.SpanFilter) {
	sp.spanFilter = newFilter
}

func (sp *AdaptiveSamplingSpanProcessor) updateTraceGraph(span *model.Span) {
	operation := &model2.Operation{
		Service: span.GetProcess().GetServiceName(),
		Name:    span.GetOperationName(),
	}
	sp.store.Add(operation)
	sp.lru.Put(span.SpanID, operation)

	parentID := span.ParentSpanID()
	if parentID != 0 {
		// only those spans which have parent spanIDs would be pushed into retryQueue for updating trace graph.
		sp.retryQueue.Produce(&retryQueueItem{
			span:     span,
			parentID: parentID,
			since:    time.Now(),
		})
	} else {
		sp.store.AddAsRoot(operation)
		sp.logger.Debug("new root operation",
			zap.String("service", operation.Service),
			zap.String("operation", operation.Name))
	}
}

func (sp *AdaptiveSamplingSpanProcessor) filterSpanToPromote(span *model.Span) {
	if sp.spanFilter.Filter(span) {
		sp.store.Promote(span)
		sp.logger.Debug("promoted operation", zap.Stringer("span", span))
	}
}

func (sp *AdaptiveSamplingSpanProcessor) processItemFromRetryQueue(item *retryQueueItem) {
	var to, from *model2.Operation

	if lruItem := sp.lru.Get(item.parentID); lruItem != nil {
		to = lruItem.(*model2.Operation)
	}
	if lruItem := sp.lru.Get(item.span.SpanID); lruItem != nil {
		from = lruItem.(*model2.Operation)
	}

	if from != nil && to != nil {
		if err := sp.store.AddEdge(from, to); err != nil {
			sp.logger.Debug("failed to add edge",
				zap.Stringer("from", from),
				zap.Stringer("to", to),
				zap.Error(err))
		}
	} else {
		if time.Now().Add(sp.retryQueueItemExpire).After(time.Now()) {
			sp.retryQueue.Produce(item)
		}
	}
}

func (sp *AdaptiveSamplingSpanProcessor) ProcessSpans(mSpans []*model.Span, options processor.SpansOptions) ([]bool, error) {
	sp.preProcessSpans(mSpans)
	sp.metrics.BatchSize.Update(int64(len(mSpans)))
	retMe := make([]bool, len(mSpans))
	for i, mSpan := range mSpans {
		ok := sp.enqueueSpan(mSpan, options.SpanFormat, options.InboundTransport)
		if !ok && sp.reportBusy {
			return nil, processor.ErrBusy
		}
		retMe[i] = ok
	}
	return retMe, nil
}

func (sp *AdaptiveSamplingSpanProcessor) Close() error {
	close(sp.stopCh)
	sp.queue.Stop()
	sp.retryQueue.Stop()
	return nil
}
