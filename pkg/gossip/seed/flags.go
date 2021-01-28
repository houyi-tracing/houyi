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

package seed

import (
	"flag"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
)

const (
	seedGrpcPort     = "gossip.seed.grpc.port"
	registryAddr     = "gossip.registry.addr"
	registryGrpcPort = "gossip.registry.grpc.port"
	lruSize          = "gossip.seed.lru.size"
)

const (
	DefaultSeedGrpcPort     = ports.SeedGrpcListenPort
	DefaultRegistryAddr     = "registry"
	DefaultRegistryGrpcPort = ports.RegistryGrpcListenPort
	DefaultLruSize          = 10000
)

type Flags struct {
	SeedGrpcPort     int
	RegistryAddress  string
	RegistryGrpcPort int
	LruSize          int
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(seedGrpcPort, DefaultSeedGrpcPort, "Port to server gRPC of local gossip seed.")
	flags.String(registryAddr, DefaultRegistryAddr, "IP or domain name of gossip registry.")
	flags.Int(registryGrpcPort, DefaultRegistryGrpcPort, "Port to server gRPC of remote gossip registry.")
	flags.Int(lruSize, DefaultLruSize, "Size of LRU for gossip message caching.")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.SeedGrpcPort = v.GetInt(seedGrpcPort)
	f.RegistryAddress = v.GetString(registryAddr)
	f.RegistryGrpcPort = v.GetInt(registryGrpcPort)
	f.LruSize = v.GetInt(lruSize)
	return f
}
