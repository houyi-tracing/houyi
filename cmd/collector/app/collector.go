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
	"github.com/houyi-tracing/houyi/cmd/collector/app/processor"
	"github.com/houyi-tracing/houyi/cmd/collector/app/server"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CollectorParams struct {
	Logger         *zap.Logger
	GrpcListenPort int
	SpanProcessor  processor.SpanProcessor
	Evaluator      evaluator.Evaluator
}

type Collector struct {
	logger         *zap.Logger
	grpcServer     *grpc.Server
	grpcListenPort int
	spanProcessor  processor.SpanProcessor
	evaluator      evaluator.Evaluator
}

func NewCollector(params *CollectorParams) *Collector {
	return &Collector{
		logger:         params.Logger,
		spanProcessor:  params.SpanProcessor,
		grpcListenPort: params.GrpcListenPort,
		evaluator:      params.Evaluator,
	}
}

func (c *Collector) Start() error {
	if gS, err := server.StartGrpcServer(&server.GrpcServerParams{
		Logger:        c.logger,
		ListenPort:    c.grpcListenPort,
		SpanProcessor: c.spanProcessor,
		Evaluator:     c.evaluator,
	}); err != nil {
		return err
	} else {
		c.grpcServer = gS
	}

	return nil
}

func (c *Collector) Close() error {
	c.grpcServer.GracefulStop()
	return nil
}
