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
	"time"
)

const (
	grpcListenPort        = "grpc.listen.port"
	httpListenPort        = "http.port"
	DefaultGrpcListenPort = ports.ConfigServerGrpcListenPort
	DefaultHttpListenPort = ports.ConfigServerHttpListenPort

	operationExpire        = "sampling.operation.expire"
	scaleFactor            = "sampling.scale.factor"
	minSamplingRate        = "sampling.min.sampling.rate"
	DefaultOperationExpire = time.Minute
	DefaultScaleFactor     = 1.0
	DefaultMinSamplingRate = 0.01

	randomPick               = "gossip.random.pick"
	probToR                  = "gossip.prob.to.r"
	heartbeatInterval        = "gossip.refresh.interval"
	DefaultRandomPick        = 5
	DefaultProbToR           = 0.25
	DefaultHeartbeatInterval = time.Second * 5
)

type Flags struct {
	GrpcListenPort    int
	HttpListenPort    int
	OperationExpire   time.Duration
	ScaleFactor       float64
	MinSamplingRate   float64
	RandomPick        int
	ProbToR           float64
	HeartbeatInterval time.Duration
}

func AddFlags(flags *flag.FlagSet) {
	flags.Int(grpcListenPort, DefaultGrpcListenPort,
		"Port to serve gRPC.")
	flags.Int(httpListenPort, DefaultHttpListenPort,
		"Port to serve HTTP.")

	flags.Duration(operationExpire, DefaultOperationExpire,
		"[Sampling] Interval for removing expired operations.")
	flags.Float64(scaleFactor, DefaultScaleFactor,
		"[Sampling] Factor used to scale sampling rates for dynamic and adaptive sampling.")
	flags.Float64(minSamplingRate, DefaultMinSamplingRate,
		"[Sampling] Minimum sampling rate for dynamic and adaptive sampling.")

	flags.Int(randomPick, DefaultRandomPick,
		"[Gossip] Number of peers a seed node send messages to when it received a message in single cycle.")
	flags.Float64(probToR, DefaultProbToR,
		"[Gossip] Probability for seed node to switch to state R when it received messages from peers.")
	flags.Duration(heartbeatInterval, DefaultHeartbeatInterval,
		"[Gossip] Interval for seed node to remove dead nodes.")
}

func (f *Flags) InitFromViper(v *viper.Viper) *Flags {
	f.GrpcListenPort = v.GetInt(grpcListenPort)
	f.HttpListenPort = v.GetInt(httpListenPort)

	f.OperationExpire = v.GetDuration(operationExpire)
	f.ScaleFactor = v.GetFloat64(scaleFactor)
	f.MinSamplingRate = v.GetFloat64(minSamplingRate)

	f.RandomPick = v.GetInt(randomPick)
	f.ProbToR = v.GetFloat64(probToR)
	f.HeartbeatInterval = v.GetDuration(heartbeatInterval)
	return f
}
