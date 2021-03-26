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

package app

import (
	"flag"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
)

const (
	collectorAddr    = "collector.addr"
	collectorPort    = "collector.port"
	configServerAddr = "config.server.addr"
	configServerPort = "config.server.port"
	grpcListenPort   = "grpc.listen.port"

	DefaultCollectorAddr    = "collector"
	DefaultCollectorPort    = ports.CollectorGrpcListenPort
	DefaultConfigServerAddr = "config-server"
	DefaultConfigServerPort = ports.ConfigServerGrpcListenPort
	DefaultGrpcListenPort   = ports.AgentGrpcListenPort
)

type Flags struct {
	CollectorAddr    string
	CollectorPort    int
	ConfigServerAddr string
	ConfigServerPort int
	GrpcListenPort   int
}

func AddFlags(flags *flag.FlagSet) {
	flags.String(collectorAddr, DefaultCollectorAddr, "IP or domain name of houyi collector.")
	flags.Int(collectorPort, DefaultCollectorPort, "Port to serve gRPC of houyi collector.")
	flags.String(configServerAddr, DefaultConfigServerAddr, "IP or domain name of configuration server.")
	flags.Int(configServerPort, DefaultConfigServerPort, "Port to serve gRPC of configuration server.")
	flags.Int(grpcListenPort, DefaultGrpcListenPort, "Port to serve gRPC of agent.")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.CollectorAddr = v.GetString(collectorAddr)
	f.CollectorPort = v.GetInt(collectorPort)
	f.ConfigServerAddr = v.GetString(configServerAddr)
	f.ConfigServerPort = v.GetInt(configServerPort)
	f.GrpcListenPort = v.GetInt(grpcListenPort)

	return f
}
