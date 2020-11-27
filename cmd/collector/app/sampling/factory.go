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
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Factory struct {
	options        Options
	logger         *zap.Logger
	metricsFactory metrics.Factory
}

func NewFactory() *Factory {
	return &Factory{
		logger:         zap.NewNop(),
		metricsFactory: metrics.NullFactory,
	}
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	AddFlags(flagSet)
}

func (f *Factory) InitFromViper(v *viper.Viper) {
	f.options = Options{}.InitFromViper(v)
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.logger = logger
	f.metricsFactory = metricsFactory
	return nil
}

func (f *Factory) CreateStrategyStore() (AdaptiveStrategyStore, error) {
	return NewAdaptiveStrategyStore(f.options, f.logger), nil
}
