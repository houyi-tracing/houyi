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

package sampling

import (
	"context"
	"fmt"
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/houyi-tracing/houyi/pkg/ds/graph"
	"github.com/houyi-tracing/houyi/pkg/ds/sst"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/thrift-gen/sampling"
	"go.uber.org/zap"
	"math"
	"sync"
	"time"
)

type qpsEntry struct {
	qps     float64
	upSince time.Time
}

// adaptiveStrategyStore indirectly stores sampling probability of services in the tree structure.
type adaptiveStrategyStore struct {
	Options

	mux                      sync.RWMutex
	logger                   *zap.Logger
	sst                      sst.SampleStrategyTree
	traceGraph               graph.ExecutionGraph
	qps                      map[string]map[string]*qpsEntry
	maxRemoteRefreshInterval time.Duration
}

func NewAdaptiveStrategyStore(opts Options, logger *zap.Logger) AdaptiveStrategyStore {
	retMe := &adaptiveStrategyStore{
		mux:                      sync.RWMutex{},
		Options:                  opts,
		logger:                   logger,
		maxRemoteRefreshInterval: opts.SamplingRefreshInterval,
		traceGraph:               graph.NewExecutionGraph(logger, opts.OperationDuration),
		sst:                      sst.NewSampleStrategyTree(opts.MaxNumChildNodes, logger),
		qps:                      make(map[string]map[string]*qpsEntry),
	}
	if retMe.MaxSamplingProbability > 1.0 {
		logger.Error("invalid maximum sampling probability (has been set to 1.0)",
			zap.Float64("invalid maximum sampling probability", retMe.MaxSamplingProbability),
			zap.Float64("new valid maximum sampling probability", 1.0))
		retMe.MaxSamplingProbability = 1.0
	}
	return retMe
}

func (ass *adaptiveStrategyStore) Add(operation *model2.Operation) {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	svc, op := operation.Service, operation.Name
	ass.traceGraph.Add(svc, op)
}

func (ass *adaptiveStrategyStore) AddAsRoot(operation *model2.Operation) {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	svc, op := operation.Service, operation.Name
	ass.sst.Add(svc, op)
	node := ass.traceGraph.Add(svc, op)
	ass.traceGraph.AddRoot(node)
}

func (ass *adaptiveStrategyStore) AddEdge(from, to *model2.Operation) error {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	if ass.traceGraph.Has(from.Service, from.Name) && ass.traceGraph.Has(to.Service, to.Name) {
		fromNode, _ := ass.traceGraph.GetNode(from.Service, from.Name)
		toNode, _ := ass.traceGraph.GetNode(to.Service, to.Name)
		ass.traceGraph.AddEdge(fromNode, toNode)
		return nil
	} else {
		return fmt.Errorf("add edge for operation not in trace graph")
	}
}

func (ass *adaptiveStrategyStore) GetRoots(operation *model2.Operation) ([]*model2.Operation, error) {
	ass.mux.RLock()
	defer ass.mux.RUnlock()

	ret := make([]*model2.Operation, 0)
	svc, op := operation.Service, operation.Name
	if node, err := ass.traceGraph.GetNode(svc, op); err == nil {
		roots := ass.traceGraph.GetRootsOf(node)
		for _, r := range roots {
			ret = append(ret, &model2.Operation{
				Service: r.Service(),
				Name:    r.Operation(),
			})
		}
		return ret, nil
	} else {
		return ret, err
	}
}

func (ass *adaptiveStrategyStore) Remove(operation *model2.Operation) error {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	svc, op := operation.Service, operation.Name
	ass.sst.Remove(svc, op)

	if node, err := ass.traceGraph.GetNode(svc, op); err == nil {
		ass.traceGraph.Remove(node)
		return nil
	} else {
		return err
	}
}

func (ass *adaptiveStrategyStore) RemoveExpired() {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	rmNodes := ass.traceGraph.RemoveExpired()
	for _, node := range rmNodes {
		ass.sst.Remove(node.Service(), node.Operation())
		ass.logger.Debug("remove expired operation", zap.Stringer("operation", node))
	}
}

func (ass *adaptiveStrategyStore) GetSamplingStrategy(_ context.Context, _ string) (*sampling.SamplingStrategyResponse, error) {
	// deprecated
	resp := &sampling.SamplingStrategyResponse{
		StrategyType:          sampling.SamplingStrategyType_PROBABILISTIC,
		ProbabilisticSampling: &sampling.ProbabilisticSamplingStrategy{SamplingRate: DefaultMinSamplingProbability},
	}
	return resp, nil
}

func (ass *adaptiveStrategyStore) GetSamplingStrategies(
	service string,
	operations model2.Operations,
	refreshInterval time.Duration) (*sampling.SamplingStrategyResponse, error) {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	// before generate strategies, we must update QPS of root operations.
	for _, op := range operations.Operations {
		if ass.traceGraph.IsRoot(service, op.Name) {
			ass.updateQpsAndRefreshInterval(service, op.Name, op.Qps, refreshInterval)
		}
	}

	strategies := make([]*sampling.OperationSamplingStrategy, 0, len(operations.Operations))
	for _, op := range operations.Operations {
		if ass.sst.Has(op.Service, op.Name) {
			if sr, err := ass.sst.GetOperationSamplingRate(op.Service, op.Name); err == nil {
				strategies = append(strategies,
					ass.newOperationSamplingStrategy(service, op.Name, sr))
				ass.logger.Debug("get sampling rate",
					zap.String("service", op.Service),
					zap.String("operation", op.Name),
					zap.Float64("QPS", op.Qps),
					zap.Float64("sampling rate", sr))
			} else {
				ass.logger.Error("failed to get sampling rate for alive operation",
					zap.Stringer("operation", op))
			}
		} else {
			strategies = append(strategies, ass.newOperationSamplingStrategy(service, op.Name, ass.MinSamplingProbability))
			ass.logger.Debug("operation is not root", zap.String("operation", op.String()))
		}

		if node, err := ass.traceGraph.GetNode(service, op.Name); err == nil {
			ass.logger.Debug("refresh root operation", zap.Stringer("operation", node))
		}
	}

	return &sampling.SamplingStrategyResponse{
		StrategyType: sampling.SamplingStrategyType_PROBABILISTIC,
		OperationSampling: &sampling.PerOperationSamplingStrategies{
			DefaultSamplingProbability: ass.MinSamplingProbability,
			PerOperationStrategies:     strategies,
		},
	}, nil
}

func (ass *adaptiveStrategyStore) Promote(span *model.Span) {
	ass.mux.Lock()
	defer ass.mux.Unlock()

	// promote the sampling rates of all roots of the node relate to inputted span.
	svc, op := span.GetProcess().GetServiceName(), span.GetOperationName()
	if node, err := ass.traceGraph.GetNode(svc, op); err == nil {
		roots := ass.traceGraph.GetRootsOf(node)
		for _, r := range roots {
			_ = ass.sst.Promote(r.Service(), r.Operation())
			ass.logger.Debug("Promoted root operation",
				zap.String("service", r.Service()),
				zap.String("operation", r.Operation()),
				zap.String("span_service", svc),
				zap.String("span_operation", op))
		}
	}
}

// updateQpsAndRefreshInterval updates the QPS and sampling refresh interval of inputted service.
func (ass *adaptiveStrategyStore) updateQpsAndRefreshInterval(service, operation string, qps float64, interval time.Duration) {
	now := time.Now()
	_, hasSvc := ass.qps[service]
	if !hasSvc {
		ass.qps[service] = make(map[string]*qpsEntry)
	}
	if qE, hasOp := ass.qps[service][operation]; hasOp {
		qE.qps = qps
		qE.upSince = now
	} else {
		ass.qps[service][operation] = &qpsEntry{
			qps:     qps,
			upSince: now,
		}
	}
	ass.maxRemoteRefreshInterval = maxDuration(ass.maxRemoteRefreshInterval, interval+time.Minute)
}

// qpsWeightCoefficient returns the weight Coefficient of service.
// The service with higher QPS would get lower qpsWeightCoefficient.
func (ass *adaptiveStrategyStore) qpsWeightCoefficient(service, operation string) float64 {
	if ass.traceGraph.IsRoot(service, operation) {
		if _, has := ass.qps[service]; has {
			if qE, has := ass.qps[service][operation]; has {
				if qE.qps == 0 {
					return 1.0
				} else {
					sum := 0.0
					for _, opMap := range ass.qps {
						for _, qps := range opMap {
							if qps.qps != 0 {
								sum += 1.0 / qps.qps
							}
						}
					}
					return (1 / qE.qps) / sum
				}
			}
		}
		ass.logger.Error("try to calculate qps weight for non-exist operation", zap.String("operation name", operation))
		return 0
	} else {
		return 1.0
	}
}

func maxDuration(d1, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	} else {
		return d2
	}
}

func (ass *adaptiveStrategyStore) newOperationSamplingStrategy(service, operation string, samplingRate float64) *sampling.OperationSamplingStrategy {
	qpsWeight := ass.qpsWeightCoefficient(service, operation)
	ret := &sampling.OperationSamplingStrategy{
		Operation: operation,
		ProbabilisticSampling: &sampling.ProbabilisticSamplingStrategy{
			SamplingRate: math.Max(math.Min(samplingRate*qpsWeight*ass.AmplificationFactor, ass.MaxSamplingProbability), ass.MinSamplingProbability),
		},
	}
	return ret
}
