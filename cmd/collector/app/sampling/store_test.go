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
	"fmt"
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/jaegertracing/jaeger/model"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ass := NewAdaptiveStrategyStore(Options{
		MaxNumChildNodes:        10,
		MaxSamplingProbability:  1,
		MinSamplingProbability:  0.01,
		TreeRefreshInterval:     time.Second * 1,
		SamplingRefreshInterval: time.Minute * 5,
	}, &zap.Logger{})

	ass.RemoveExpired()
}

func TestPromote(t *testing.T) {
	logger, _ := zap.NewProduction()
	ass := NewAdaptiveStrategyStore(Options{
		MaxNumChildNodes:        4,
		MaxSamplingProbability:  1,
		MinSamplingProbability:  0.01,
		TreeRefreshInterval:     time.Second * 1,
		SamplingRefreshInterval: time.Minute * 5,
		AmplificationFactor:     1,
	}, logger)

	ops := make([]model2.Operation, 0)

	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		ops = append(ops, model2.Operation{
			Service: fmt.Sprintf("svc"),
			Name:    fmt.Sprintf("op%d", n),
			Qps:     10,
		})
	}

	for _, op := range ops {
		ass.AddAsRoot(&op)
		_, _ = ass.GetSamplingStrategies("svc", model2.Operations{Operations: ops}, time.Minute)
	}

	sr, _ := ass.GetSamplingStrategies("svc", model2.Operations{
		Operations: []model2.Operation{
			{
				"svc", "op1", 1,
			},
		},
	}, time.Minute)
	fmt.Println(sr)

	ass.Promote(&model.Span{
		OperationName: "op1",
		Process: &model.Process{
			ServiceName: "svc",
		},
	})
	sr, _ = ass.GetSamplingStrategies("svc", model2.Operations{
		Operations: []model2.Operation{
			{
				"svc", "op1", 1,
			},
		},
	}, time.Minute)
	fmt.Println(sr)
}
