// Copyright (c) 2020 The Jaeger Authors.
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
	"context"
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/jaegertracing/jaeger/model"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics/fork"
	"github.com/uber/jaeger-lib/metrics/metricstest"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/pkg/healthcheck"
	"github.com/jaegertracing/jaeger/thrift-gen/sampling"
)

var _ (io.Closer) = (*Collector)(nil)

func TestNewCollector(t *testing.T) {
	// prepare
	hc := healthcheck.New()
	logger := zap.NewNop()
	baseMetrics := metricstest.NewFactory(time.Hour)
	spanWriter := &fakeSpanWriter{}
	strategyStore := &mockStrategyStore{}

	c := New(&CollectorParams{
		ServiceName:    "collector",
		Logger:         logger,
		MetricsFactory: baseMetrics,
		SpanWriter:     spanWriter,
		StrategyStore:  strategyStore,
		HealthCheck:    hc,
	})
	collectorOpts := &CollectorOptions{}

	// test
	_ = c.Start(collectorOpts)

	// verify
	assert.NoError(t, c.Close())
}

type mockStrategyStore struct {
}

func (m *mockStrategyStore) Add(_ *model2.Operation) {
	panic("implement me")
}

func (m *mockStrategyStore) AddEdge(_ *model2.Operation, _ *model2.Operation) error {
	panic("implement me")
}

func (m *mockStrategyStore) AddAsRoot(_ *model2.Operation) {
	panic("implement me")
}

func (m *mockStrategyStore) GetRoots(_ *model2.Operation) ([]*model2.Operation, error) {
	panic("implement me")
}

func (m *mockStrategyStore) Remove(_ *model2.Operation) error {
	panic("implement me")
}

func (m *mockStrategyStore) RemoveExpired() {
	panic("implement me")
}

func (m *mockStrategyStore) GetSamplingStrategies(_ string, _ model2.Operations, _ time.Duration) (*sampling.SamplingStrategyResponse, error) {
	return &sampling.SamplingStrategyResponse{}, nil
}

func (m *mockStrategyStore) Promote(_ *model.Span) {

}

func (m *mockStrategyStore) Start() error {
	return nil
}

func (m *mockStrategyStore) Close() error {
	return nil
}

func (m *mockStrategyStore) GetSamplingStrategy(_ context.Context, _ string) (*sampling.SamplingStrategyResponse, error) {
	return &sampling.SamplingStrategyResponse{}, nil
}

func TestCollector_PublishOpts(t *testing.T) {
	// prepare
	hc := healthcheck.New()
	logger := zap.NewNop()
	baseMetrics := metricstest.NewFactory(time.Second)
	forkFactory := metricstest.NewFactory(time.Second)
	metricsFactory := fork.New("internal", forkFactory, baseMetrics)
	spanWriter := &fakeSpanWriter{}
	strategyStore := &mockStrategyStore{}

	c := New(&CollectorParams{
		ServiceName:    "collector",
		Logger:         logger,
		MetricsFactory: metricsFactory,
		SpanWriter:     spanWriter,
		StrategyStore:  strategyStore,
		HealthCheck:    hc,
	})
	collectorOpts := &CollectorOptions{
		NumWorkers: 24,
		QueueSize:  42,
	}

	_ = c.Start(collectorOpts)

	forkFactory.AssertGaugeMetrics(t, metricstest.ExpectedMetric{
		Name:  "internal.collector.num-workers",
		Value: 24,
	})
	forkFactory.AssertGaugeMetrics(t, metricstest.ExpectedMetric{
		Name:  "internal.collector.queue-size",
		Value: 42,
	})

	_ = c.Close()
}
