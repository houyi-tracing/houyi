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

package handler

import (
	"context"
	"github.com/houyi-tracing/houyi/cmd/agent/app/transport"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	jaeger "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	api_v1.UnimplementedStrategyManagerServer

	logger                   *zap.Logger
	collectorTransport       *transport.CollectorTransport
	strategyManagerTransport *transport.StrategyManagerTransport
}

func NewGrpcHandler(logger *zap.Logger, ct *transport.CollectorTransport, smt *transport.StrategyManagerTransport) *GrpcHandler {
	return &GrpcHandler{
		logger:                   logger,
		collectorTransport:       ct,
		strategyManagerTransport: smt,
	}
}

func (p *GrpcHandler) GetStrategy(ctx context.Context, request *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	p.logger.Debug("GetStrategy", zap.String("request", request.String()))

	resp, err := p.strategyManagerTransport.GetStrategy(ctx, request)
	if err != nil {
		p.logger.Error("failed to get strategies", zap.Error(err))
	}
	return resp, err
}

func (p *GrpcHandler) PostSpans(ctx context.Context, request *jaeger.PostSpansRequest) (*jaeger.PostSpansResponse, error) {
	resp, err := p.collectorTransport.PostSpans(ctx, request)
	if err != nil {
		p.logger.Error("failed to post spans", zap.Error(err))
	}
	return resp, err
}
