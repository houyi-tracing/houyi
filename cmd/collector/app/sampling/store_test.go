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
