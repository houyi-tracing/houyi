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
	"flag"
	"github.com/spf13/viper"
	"time"
)

const (
	maxNumChildNodes                     = "sampling.max-num-child-nodes"
	minSamplingProbability               = "sampling.max-samp-prob"
	maxSamplingProbability               = "sampling.min-samp-prob"
	treeRefreshInterval                  = "sampling.tree-refresh-interval"
	samplingRefreshInterval              = "sampling.sampling-refresh-interval"
	amplificationFactor                  = "sampling.amplification-factor"
	samplingRefreshIntervalShrinkageRate = "sampling.shrinkage-rate"
	operationDuration                    = "operation.duration"

	DefaultMaxNumChildNodes                     = 4
	DefaultMinSamplingProbability               = 0.001
	DefaultMaxSamplingProbability               = 1.0
	DefaultTreeRefreshInterval                  = time.Minute * 2
	DefaultSamplingRefreshInterval              = time.Minute * 1
	DefaultAmplificationFactor                  = 1.0
	DefaultSamplingRefreshIntervalShrinkageRate = 0.95
	DefaultOperationDuration                    = time.Minute
)

type Options struct {
	AmplificationFactor                  float64
	MaxNumChildNodes                     int
	MaxSamplingProbability               float64
	MinSamplingProbability               float64
	TreeRefreshInterval                  time.Duration
	SamplingRefreshInterval              time.Duration
	SamplingRefreshIntervalShrinkageRate float64
	OperationDuration                    time.Duration
}

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.Float64(amplificationFactor, DefaultAmplificationFactor,
		"Amplification factor to amplify sampling rate")
	flagSet.Int(maxNumChildNodes, DefaultMaxNumChildNodes,
		"Maximum number of child nodes for initializing sample strategy tree.",
	)
	flagSet.Float64(minSamplingProbability, DefaultMinSamplingProbability,
		"Minimum sampling probability for all operations.",
	)
	flagSet.Float64(maxSamplingProbability, DefaultMaxSamplingProbability,
		"Maximum sampling probability for all operations.",
	)
	flagSet.Duration(treeRefreshInterval, DefaultTreeRefreshInterval,
		"Refresh interval of sample adaptiveStrategyStore tree.",
	)
	flagSet.Duration(samplingRefreshInterval, DefaultSamplingRefreshInterval,
		"Refresh interval between different requests to get sampling strategy for same service.",
	)
	flagSet.Float64(samplingRefreshIntervalShrinkageRate, DefaultSamplingRefreshIntervalShrinkageRate,
		"Shrinkage rate for reducing maximum sampling refresh interval in each time of prune operation.")
	flagSet.Duration(operationDuration, DefaultOperationDuration,
		"Operation duration for removing expired operations in trace graph.")
}

func (opts Options) InitFromViper(v *viper.Viper) Options {
	opts.MaxNumChildNodes = v.GetInt(maxNumChildNodes)
	opts.MinSamplingProbability = v.GetFloat64(minSamplingProbability)
	opts.MaxSamplingProbability = v.GetFloat64(maxSamplingProbability)
	opts.TreeRefreshInterval = v.GetDuration(treeRefreshInterval)
	opts.SamplingRefreshInterval = v.GetDuration(samplingRefreshInterval)
	opts.AmplificationFactor = v.GetFloat64(amplificationFactor)
	opts.SamplingRefreshIntervalShrinkageRate = v.GetFloat64(samplingRefreshIntervalShrinkageRate)
	opts.OperationDuration = v.GetDuration(operationDuration)
	return opts
}
