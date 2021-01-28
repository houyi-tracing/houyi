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

package app

import (
	"flag"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
	"time"
)

const (
	randomPick      = "random.pick"
	probToR         = "prob.to.R"
	refreshInterval = "refresh.interval"
	listenPort      = "listen.port"
)

const (
	DefaultRandomPick      = 5
	DefaultProbToR         = 0.25
	DefaultRefreshInterval = time.Second * 5
	DefaultListenPort      = ports.RegistryGrpcListenPort
)

type Flags struct {
	RandomPick      int
	ListenPort      int
	ProbToR         float64
	RefreshInterval time.Duration
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(randomPick, DefaultRandomPick,
		"Number of peers a seed node send messages to when it received a message in single cycle.")
	flags.Int(listenPort, DefaultListenPort,
		"Port to server gRPC of registry")
	flags.Float64(probToR, DefaultProbToR,
		"Probability for seed node to switch to state R when it received messages from peers.")
	flags.Duration(refreshInterval, DefaultRefreshInterval,
		"Interval for seed node to remove dead nodes.")
}

func (rOpts *Flags) InitFromViper(v *viper.Viper) *Flags {
	rOpts.RandomPick = v.GetInt(randomPick)
	rOpts.ListenPort = v.GetInt(listenPort)
	rOpts.ProbToR = v.GetFloat64(probToR)
	rOpts.RefreshInterval = v.GetDuration(refreshInterval)
	return rOpts
}
