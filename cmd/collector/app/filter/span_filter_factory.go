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

package filter

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type spanFilterFactory struct {
	configFile string
	logger     *zap.Logger
}

func NewFilterFactory() Factory {
	return &spanFilterFactory{}
}

func (s *spanFilterFactory) Initialize(v *viper.Viper, logger *zap.Logger) {
	if logger == nil {
		s.logger = zap.NewNop()
	} else {
		s.logger = logger
	}

	opts := new(options).InitFromViper(v)
	if opts.FilterConfigFile == "" {
		s.logger.Info("filter tags file is empty")
	} else {
		s.configFile = opts.FilterConfigFile
	}
}

func (s *spanFilterFactory) CreateFilter() SpanFilter {
	if s.configFile == "" {
		return NewNullSpanFilter()
	}
	return newTagSpanFilterFromFile(s.configFile, s.logger)
}
