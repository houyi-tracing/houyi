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
	"flag"
	"github.com/spf13/viper"
)

const (
	filterConfigFile = "filter.config.file"

	defaultConfigFile = "${HOME}/filter-config.json"
)

type options struct {
	FilterConfigFile string
}

func AddFlags(flags *flag.FlagSet) {
	flags.String(
		filterConfigFile,
		defaultConfigFile,
		"Configuration file for span filter to filter span for promoting operation in SST.")
}

func (opts options) InitFromViper(v *viper.Viper) options {
	opts.FilterConfigFile = v.GetString(filterConfigFile)
	return opts
}
