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
	"context"
	"fmt"
	"github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/queue"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	ParentTagNameService   = "p-svc"
	ParentTagNameOperation = "p-op"
)

type queueItem struct {
	queuedTime time.Time
	span       *model.Span
}

type spanProcessor struct {
	logger *zap.Logger

	workers int

	strategyManagerEndpoint *routing.Endpoint

	queue queue.DynamicQueue

	filterSpan   filter.FilterSpan
	evaluateSpan evaluator.EvaluateSpan
	processSpan  ProcessSpan
	spanWriter   spanstore.Writer

	traceGraph tg.TraceGraph
	seed       gossip.Seed

	stopCh chan *sync.WaitGroup
	opCh   chan *api_v1.Operation
}

func NewSpanProcessor(logger *zap.Logger, opts ...Option) SpanProcessor {
	sp := newSpanProcessor(logger, opts...)

	sp.queue.StartConsumers(sp.workers, func(item interface{}) {
		if i, ok := item.(*queueItem); ok {
			sp.processItemFromQueue(i)
		}
	})

	go sp.strategyManagerClient()

	return sp
}

func newSpanProcessor(logger *zap.Logger, opts ...Option) *spanProcessor {
	o := new(options).apply(opts...)
	sp := &spanProcessor{
		logger:                  logger,
		filterSpan:              o.filterSpan,
		evaluateSpan:            o.evaluateSpan,
		spanWriter:              o.spanWriter,
		strategyManagerEndpoint: o.strategyManagerEndpoint,
		queue:                   queue.NewDynamicQueue(),
		traceGraph:              o.traceGraph,
		seed:                    o.seed,
		opCh:                    make(chan *api_v1.Operation, 1000),
		stopCh:                  make(chan *sync.WaitGroup),
	}
	processSpanFuncs := []ProcessSpan{sp.parseSpan, sp.saveSpan}
	sp.processSpan = ChainedProcessSpan(processSpanFuncs...)

	return sp
}

func (sp *spanProcessor) ProcessSpans(spans []*model.Span) error {
	for _, span := range spans {
		if ok := sp.enqueueSpan(span); !ok {
			return fmt.Errorf("span processor is busy")
		}
	}
	return nil
}

func (sp *spanProcessor) Close() error {
	sp.queue.Stop()

	var wg sync.WaitGroup
	wg.Add(1)
	// stop strategy manager client
	sp.stopCh <- &wg
	wg.Wait()

	if err := sp.seed.Stop(); err != nil {
		return err
	}

	return nil
}

func (sp *spanProcessor) saveSpan(span *model.Span) {
	if nil == span.Process {
		sp.logger.Error("process is empty for the span")
		return
	}

	if err := sp.spanWriter.WriteSpan(context.TODO(), span); err != nil {
		sp.logger.Error("Failed to save span", zap.Error(err))
	}
}

func (sp *spanProcessor) enqueueSpan(span *model.Span) bool {
	if sp.filterSpan(span) {
		// as in "not dropped", because it's actively rejected. [Jaeger]
		return true
	}

	item := &queueItem{
		queuedTime: time.Now(),
		span:       span,
	}
	return sp.queue.Produce(item)
}

func (sp *spanProcessor) parseSpan(span *model.Span) {
	pSvc, pOp := getTagStrVal(span, ParentTagNameService), getTagStrVal(span, ParentTagNameOperation)

	if pSvc == "" || pOp == "" {
		return
	}

	parentOp := &api_v1.Operation{
		Service:   pSvc,
		Operation: pOp,
	}
	currOp := &api_v1.Operation{
		Service:   span.GetOperationName(),
		Operation: span.GetProcess().ServiceName,
	}
	rel := &api_v1.Relation{
		From: parentOp,
		To:   currOp,
	}

	if !sp.traceGraph.HasRelation(rel) {
		sp.seed.MongerNewRelation(rel)
	}

	// Evaluate a span whether it is need to be promoted
	if sp.evaluateSpan(span) {
		sp.opCh <- currOp
	}
}

func (sp *spanProcessor) strategyManagerClient() {
	for {
		select {
		case op := <-sp.opCh:
			sp.promoteOperation(op)
		case wg := <-sp.stopCh:
			wg.Done()
			return
		}
	}
}

func (sp *spanProcessor) promoteOperation(op *api_v1.Operation) {
	conn, err := grpc.Dial(sp.strategyManagerEndpoint.String(), grpc.WithInsecure(), grpc.WithBlock())
	if conn == nil || err != nil {
		sp.logger.Fatal("Could not dial to strategy manager",
			zap.String("address", sp.strategyManagerEndpoint.String()))
	} else {
		defer conn.Close()
	}

	c := api_v1.NewDynamicStrategyManagerClient(conn)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	req := op
	if _, err := c.Promote(ctx, req); err != nil {
		sp.logger.Error("Failed to send promote request to strategy manager", zap.Error(err))
	} else {
		sp.logger.Debug("Received promote reply from strategy manager",
			zap.String("operation", op.String()))
	}
}

func (sp *spanProcessor) processItemFromQueue(item *queueItem) {
	sp.processSpan(item.span)
}

func getTagStrVal(span *model.Span, tagName string) string {
	tags := span.GetTags()
	for _, t := range tags {
		if t.Key == tagName && t.VType == model.ValueType_STRING {
			return t.VStr
		}
	}
	return ""
}
