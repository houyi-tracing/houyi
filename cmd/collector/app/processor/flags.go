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

package processor

import (
	"flag"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
)

const (
	numWorkers         = "num.workers"
	strategyMangerAddr = "strategy.manager.addr"
	strategyMangerPort = "strategy.manager.port"

	DefaultNumWorkers          = 4
	DefaultStrategyManagerAddr = "strategy-manager"
	DefaultStrategyManagerPort = ports.StrategyManagerGrpcListenPort
)

type Flags struct {
	NumWorkers          int
	StrategyManagerAddr string
	StrategyManagerPort int
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(numWorkers,
		DefaultNumWorkers, "Number of workers to consume dynamic queue in span processor.")
	flags.String(strategyMangerAddr, DefaultStrategyManagerAddr, "IP or domain name of strategy manager.")
	flags.Int(strategyMangerPort, DefaultStrategyManagerPort, "Port to server gRPC for strategy manager")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.NumWorkers = v.GetInt(numWorkers)
	f.StrategyManagerAddr = v.GetString(strategyMangerAddr)
	f.StrategyManagerPort = v.GetInt(strategyMangerPort)

	return f
}
