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
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
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
	sync.RWMutex
	Options

	logger                   *zap.Logger
	sst                      sst.SampleStrategyTree
	srCache                  map[string]map[string]float64 // sampling rates cache generated from SST
	qps                      map[string]map[string]*qpsEntry
	stopCh                   chan struct{}
	maxRemoteRefreshInterval time.Duration
}

func NewAdaptiveStrategyStore(opts Options, logger *zap.Logger) AdaptiveStrategyStore {
	retMe := &adaptiveStrategyStore{
		Options:                  opts,
		logger:                   logger,
		maxRemoteRefreshInterval: opts.SamplingRefreshInterval,
		sst:                      sst.NewSampleStrategyTree(opts.MaxNumChildNodes, logger),
		srCache:                  make(map[string]map[string]float64),
		stopCh:                   make(chan struct{}),
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

func (ass *adaptiveStrategyStore) GetSamplingStrategy(_ context.Context, service string) (*sampling.SamplingStrategyResponse, error) {
	ass.Lock()
	defer ass.Unlock()

	if !ass.sst.HasService(service) {
		ass.sst.AddService(service)
	}
	strategies := make([]*sampling.OperationSamplingStrategy, 0, len(ass.srCache[service]))
	for op, sr := range ass.srCache[service] {
		strategies = append(strategies, ass.newOperationSamplingStrategy(service, op, sr))
	}
	resp := &sampling.SamplingStrategyResponse{
		StrategyType: sampling.SamplingStrategyType_PROBABILISTIC,
		OperationSampling: &sampling.PerOperationSamplingStrategies{
			DefaultSamplingProbability: ass.MinSamplingProbability,
			PerOperationStrategies:     strategies,
		},
	}

	return resp, nil
}

func (ass *adaptiveStrategyStore) GetSamplingStrategies(
	service string,
	operations model2.Operations,
	refreshInterval time.Duration) (*sampling.SamplingStrategyResponse, error) {
	ass.Lock()
	defer ass.Unlock()

	if !ass.sst.HasService(service) {
		ass.sst.AddService(service)
	}

	hasNewOperation := false
	for _, op := range operations.Operations {
		if !ass.sst.Has(service, op.Name) {
			ass.sst.Add(service, op.Name)
			if _, has := ass.qps[service]; !has {
				ass.qps[service] = make(map[string]*qpsEntry)
			}
			ass.qps[service][op.Name] = &qpsEntry{
				qps:     op.Qps,
				upSince: time.Now(),
			}
			hasNewOperation = true
			ass.logger.Debug("new op",
				zap.String("service", service),
				zap.String("op", op.Name))
		}
	}

	if hasNewOperation {
		ass.rebuildSSTOutputCache()
	}

	// before generate strategies, we must update QPS of inputted operations.
	ass.updateQpsAndRefreshInterval(service, operations.Operations, refreshInterval)

	strategies := make([]*sampling.OperationSamplingStrategy, 0, len(operations.Operations))
	for _, op := range operations.Operations {
		if sr, has := ass.srCache[service][op.Name]; has {
			strategies = append(strategies, ass.newOperationSamplingStrategy(service, op.Name, sr))
		} else {
			ass.logger.Fatal("operation not exist in sst output cache", zap.String("operation", op.Name))
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
	ass.Lock()
	defer ass.Unlock()

	serviceName := span.GetProcess().GetServiceName()
	_ = ass.sst.Promote(serviceName, span.OperationName)
	ass.rebuildSSTOutputCache()
}

func (ass *adaptiveStrategyStore) Start() error {
	ass.refresh(ass.TreeRefreshInterval)
	ass.logger.Info("Started adaptive strategy sst",
		zap.Float64("max sampling probability", ass.MaxSamplingProbability),
		zap.Float64("min sampling probability", ass.MinSamplingProbability),
		zap.Int("max number of child nodes in SST", ass.MaxNumChildNodes))
	return nil
}

func (ass *adaptiveStrategyStore) Close() error {
	close(ass.stopCh)
	return nil
}

// updateQpsAndRefreshInterval updates the QPS and sampling refresh interval of inputted service.
func (ass *adaptiveStrategyStore) updateQpsAndRefreshInterval(service string, operations []model2.Operation, interval time.Duration) {
	now := time.Now()

	for _, op := range operations {
		_, ok := ass.qps[service]
		if !ok {
			ass.qps[service] = make(map[string]*qpsEntry)
		}
		if qE, has := ass.qps[service][op.Name]; !has {
			qE.qps = op.Qps
			qE.upSince = now
		} else {
			ass.qps[service][op.Name] = &qpsEntry{
				qps:     op.Qps,
				upSince: now,
			}
		}
	}

	ass.maxRemoteRefreshInterval = maxDuration(ass.maxRemoteRefreshInterval, interval)
}

// rebuildSSTOutputCache must be called after each time of pruning, adding or promoting operations to SST.
func (ass *adaptiveStrategyStore) rebuildSSTOutputCache() {
	ass.srCache = ass.sst.GetSamplingRate()
}

func (ass *adaptiveStrategyStore) prune() {
	ass.Lock()
	defer ass.Unlock()

	ass.sst.Prune(ass.MinSamplingProbability)
	now := time.Now()

	for svc, opMap := range ass.qps {
		for op, qE := range opMap {
			if now.Sub(qE.upSince) > ass.maxRemoteRefreshInterval {
				if innerOpMap, ok := ass.srCache[svc]; ok {
					delete(innerOpMap, op)
				}
				delete(ass.qps[svc], op)
				ass.sst.Remove(svc, op)
				ass.logger.Debug("removed expired operation.",
					zap.String("service", svc),
					zap.String("operation", op))
			}
		}
	}

	ass.maxRemoteRefreshInterval = time.Duration(
		float64(ass.maxRemoteRefreshInterval.Nanoseconds()) * ass.SamplingRefreshIntervalShrinkageRate)
}

// refresh refresh response cache and QPS periodically.
func (ass *adaptiveStrategyStore) refresh(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ass.prune()
				ass.logger.Debug("report strategy sst",
					zap.Int("operations", len(ass.qps)),
					zap.Duration("max refresh interval", ass.maxRemoteRefreshInterval))
			case <-ass.stopCh:
				return
			}
		}
	}()
}

// qpsWeightCoefficient returns the weight Coefficient of service.
// The service with higher QPS would get lower qpsWeightCoefficient.
func (ass *adaptiveStrategyStore) qpsWeightCoefficient(service, operation string) float64 {
	// TODO find a proper function to calculate weight coefficient for operations
	if qE, has := ass.qps[service][operation]; has {
		return math.Exp(-qE.qps * ass.ExpCoefficient)
	} else {
		ass.logger.Fatal("try to calculate qps weight for non-exist operation", zap.String("operation name", operation))
		return 0
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
	return &sampling.OperationSamplingStrategy{
		Operation: operation,
		ProbabilisticSampling: &sampling.ProbabilisticSamplingStrategy{
			SamplingRate: math.Max(math.Min(samplingRate*ass.qpsWeightCoefficient(service, operation)*ass.AmplificationFactor, ass.MaxSamplingProbability), ass.MinSamplingProbability),
		},
	}
}
