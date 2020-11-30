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
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/jaegertracing/jaeger/cmd/collector/app/sampling/strategystore"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/thrift-gen/sampling"
	"time"
)

type AdaptiveStrategyStore interface {
	strategystore.StrategyStore

	Add(*model2.Operation)
	AddEdge(from *model2.Operation, to *model2.Operation) error
	AddAsRoot(*model2.Operation)
	GetRoots(*model2.Operation) ([]*model2.Operation, error)
	Remove(*model2.Operation) error
	RemoveExpired()

	// GetSamplingStrategies returns sampling strategies of operations of inputted service and refresh the maximum
	// sampling pulling interval via the inputted interval.
	GetSamplingStrategies(string, model2.Operations, time.Duration) (*sampling.SamplingStrategyResponse, error)

	// Promote promotes the level of operation relate to inputted span. As a result, the sampling rate of this operation
	// would be increased if the operation is not at the top level of sampling strategy tree.
	Promote(*model.Span)
}
