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
	grpcListenPort = "collector.grpc.port"

	DefaultGrpcListenPort = ports.CollectorGrpcListenPort
)

type Flags struct {
	GrpcListenPort int
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(grpcListenPort, DefaultGrpcListenPort, "Port to server gRPC of collector.")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.GrpcListenPort = v.GetInt(grpcListenPort)
	return f
}
