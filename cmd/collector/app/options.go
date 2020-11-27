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
	filter "github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/cmd/collector/app/sampling"
	"github.com/jaegertracing/jaeger/cmd/collector/app"
	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	"github.com/jaegertracing/jaeger/cmd/collector/app/sanitizer"
	"github.com/jaegertracing/jaeger/model"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"time"
)

const (
	// DefaultNumWorkers is the default number of workers consuming from the processor queue
	DefaultNumWorkers = 50
	// DefaultQueueSize is the size of the processor's queue
	DefaultQueueSize = 2000
	// DefaultLruCapacity is the capacity of LRU which contains key-value pairs
	DefaultLruCapacity = 10000
	// DefaultOperationDuration is the expiration duration of span in trace graph
	DefaultOperationDuration = time.Minute * 5
	// DefaultMaxRetries is the maximum retires of span to update trace graph
	DefaultMaxRetires = 10
	// DefaultFilterTagsFileName is the default filter tags filename for filtering spans
	DefaultFilterTagsFileName = "filter-tags.json"
	// DefaultRetryQueueNumWorkers is the default number of workers consuming from the retry queue
	DefaultRetryQueueNumWorkers = 5
)

type options struct {
	logger             *zap.Logger
	serviceMetrics     metrics.Factory
	hostMetrics        metrics.Factory
	preProcessSpans    app.ProcessSpans
	sanitizer          sanitizer.SanitizeSpan
	preSave            app.ProcessSpan
	spanFilter         app.FilterSpan
	numWorkers         int
	blockingSubmit     bool
	queueSize          int
	dynQueueSizeWarmup uint
	dynQueueSizeMemory uint
	reportBusy         bool
	extraFormatTypes   []processor.SpanFormat
	collectorTags      map[string]string

	// Adaptive sampling
	lruCapacity          int
	maxRetires           int
	operationDuration    time.Duration
	asSpanFilter         filter.SpanFilter // as, short for adaptive sampling
	filterTagsFilename   string
	store                sampling.AdaptiveStrategyStore
	retryQueueNumWorkers int
}

// Option is a function that sets some option on StorageBuilder.
type Option func(c *options)

// Options is a factory for all available Option's
var Options options

// Logger creates a Option that initializes the logger
func (options) Logger(logger *zap.Logger) Option {
	return func(b *options) {
		b.logger = logger
	}
}

// ServiceMetrics creates an Option that initializes the serviceMetrics metrics factory
func (options) ServiceMetrics(serviceMetrics metrics.Factory) Option {
	return func(b *options) {
		b.serviceMetrics = serviceMetrics
	}
}

// HostMetrics creates an Option that initializes the hostMetrics metrics factory
func (options) HostMetrics(hostMetrics metrics.Factory) Option {
	return func(b *options) {
		b.hostMetrics = hostMetrics
	}
}

// PreProcessSpans creates an Option that initializes the preProcessSpans function
func (options) PreProcessSpans(preProcessSpans app.ProcessSpans) Option {
	return func(b *options) {
		b.preProcessSpans = preProcessSpans
	}
}

// Sanitizer creates an Option that initializes the sanitizer function
func (options) Sanitizer(sanitizer sanitizer.SanitizeSpan) Option {
	return func(b *options) {
		b.sanitizer = sanitizer
	}
}

// PreSave creates an Option that initializes the preSave function
func (options) PreSave(preSave app.ProcessSpan) Option {
	return func(b *options) {
		b.preSave = preSave
	}
}

// SpanFilter creates an Option that initializes the tagSpanFilter function
func (options) SpanFilter(spanFilter app.FilterSpan) Option {
	return func(b *options) {
		b.spanFilter = spanFilter
	}
}

// NumWorkers creates an Option that initializes the number of queue consumers AKA workers
func (options) NumWorkers(numWorkers int) Option {
	return func(b *options) {
		b.numWorkers = numWorkers
	}
}

// BlockingSubmit creates an Option that initializes the blockingSubmit boolean
func (options) BlockingSubmit(blockingSubmit bool) Option {
	return func(b *options) {
		b.blockingSubmit = blockingSubmit
	}
}

// QueueSize creates an Option that initializes the queue size
func (options) QueueSize(queueSize int) Option {
	return func(b *options) {
		b.queueSize = queueSize
	}
}

// DynQueueSize creates an Option that initializes the queue size
func (options) DynQueueSizeWarmup(dynQueueSizeWarmup uint) Option {
	return func(b *options) {
		b.dynQueueSizeWarmup = dynQueueSizeWarmup
	}
}

// DynQueueSize creates an Option that initializes the queue size
func (options) DynQueueSizeMemory(dynQueueSizeMemory uint) Option {
	return func(b *options) {
		b.dynQueueSizeMemory = dynQueueSizeMemory
	}
}

// ReportBusy creates an Option that initializes the reportBusy boolean
func (options) ReportBusy(reportBusy bool) Option {
	return func(b *options) {
		b.reportBusy = reportBusy
	}
}

// ExtraFormatTypes creates an Option that initializes the extra list of format types
func (options) ExtraFormatTypes(extraFormatTypes []processor.SpanFormat) Option {
	return func(b *options) {
		b.extraFormatTypes = extraFormatTypes
	}
}

// CollectorTags creates an Option that initializes the extra tags to append to the spans flowing through this collector
func (options) CollectorTags(extraTags map[string]string) Option {
	return func(b *options) {
		b.collectorTags = extraTags
	}
}

// LruCapacity creates an Option that initializes the capacity of LRU in adaptive sampling
// span processors
func (options) LruCapacity(capacity int) Option {
	return func(o *options) {
		o.lruCapacity = capacity
	}
}

// MaxRetries creates an Option that initializes the maximum retries for update trace graph
// in adaptive sampling span processor.
func (options) MaxRetries(retries int) Option {
	return func(o *options) {
		o.maxRetires = retries
	}
}

// OperationDuration creates an Option that initializes the expire minutes of span in trace graph
// in adaptive sampling span processor.
func (options) OperationDuration(duration time.Duration) Option {
	return func(o *options) {
		o.operationDuration = duration
	}
}

// AdaptiveSamplingSpanFilter creates an Option that initializes span filter for adaptive sampling.
func (options) AdaptiveSamplingSpanFilter(filter filter.SpanFilter) Option {
	return func(o *options) {
		o.asSpanFilter = filter
	}
}

func (options) FilterTagsFilename(filename string) Option {
	return func(o *options) {
		o.filterTagsFilename = filename
	}
}

func (options) StrategyStore(store sampling.AdaptiveStrategyStore) Option {
	return func(o *options) {
		o.store = store
	}
}

func (options) RetryQueueNumWorkers(retryQueueNumWorkers int) Option {
	return func(c *options) {
		c.retryQueueNumWorkers = retryQueueNumWorkers
	}
}

func (o options) apply(opts ...Option) options {
	ret := options{}
	for _, opt := range opts {
		opt(&ret)
	}
	if ret.logger == nil {
		ret.logger = zap.NewNop()
	}
	if ret.serviceMetrics == nil {
		ret.serviceMetrics = metrics.NullFactory
	}
	if ret.hostMetrics == nil {
		ret.hostMetrics = metrics.NullFactory
	}
	if ret.preProcessSpans == nil {
		ret.preProcessSpans = func(spans []*model.Span) {}
	}
	if ret.sanitizer == nil {
		ret.sanitizer = func(span *model.Span) *model.Span { return span }
	}
	if ret.preSave == nil {
		ret.preSave = func(span *model.Span) {}
	}
	if ret.spanFilter == nil {
		ret.spanFilter = func(span *model.Span) bool { return true }
	}
	if ret.numWorkers == 0 {
		ret.numWorkers = DefaultNumWorkers
	}

	// adaptive sampling
	if ret.lruCapacity == 0 {
		ret.lruCapacity = DefaultLruCapacity
	}
	if ret.maxRetires == 0 {
		ret.maxRetires = DefaultMaxRetires
	}
	if ret.operationDuration == 0 {
		ret.operationDuration = DefaultOperationDuration
	}
	if ret.asSpanFilter == nil {
		ret.asSpanFilter = filter.NewNullSpanFilter()
	}
	if ret.filterTagsFilename == "" {
		ret.filterTagsFilename = DefaultFilterTagsFileName
	}
	if ret.store == nil {
		ret.store = nil
	}
	if ret.retryQueueNumWorkers == 0 {
		ret.retryQueueNumWorkers = DefaultRetryQueueNumWorkers
	}
	return ret
}
