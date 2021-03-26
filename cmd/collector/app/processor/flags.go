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
	numWorkers       = "num.workers"
	configServerAddr = "sampling.config.server.addr"
	configServerPort = "sampling.config.server.port"

	DefaultNumWorkers       = 4
	DefaultConfigServerAddr = "config-server"
	DefaultConfigServerPort = ports.ConfigServerGrpcListenPort
)

type Flags struct {
	NumWorkers       int
	ConfigServerAddr string
	ConfigServerPort int
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(numWorkers,
		DefaultNumWorkers, "Number of workers to consume dynamic queue in span processor.")
	flags.String(configServerAddr, DefaultConfigServerAddr, "[Sampling] IP or domain name of configuration server.")
	flags.Int(configServerPort, DefaultConfigServerPort, "[Sampling] Port to server gRPC for configuration server.")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.NumWorkers = v.GetInt(numWorkers)
	f.ConfigServerAddr = v.GetString(configServerAddr)
	f.ConfigServerPort = v.GetInt(configServerPort)

	return f
}
