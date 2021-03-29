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

package model

const (
	StrategyType_Const        = "const"
	StrategyType_Default      = "default"
	StrategyType_Probability  = "probability"
	StrategyType_RateLimiting = "rate-limiting"
	StrategyType_Adaptive     = "adaptive"
	StrategyType_Dynamic      = "dynamic"
)

type Strategy struct {
	Service            string  `json:"service"`
	Operation          string  `json:"operation"`
	Type               string  `json:"type"`
	SamplingRate       float64 `json:"sampling_rate"`
	MaxTracesPerSecond int64   `json:"max_traces_per_second"`
	AlwaysSample       bool    `json:"always_sample"`
}
