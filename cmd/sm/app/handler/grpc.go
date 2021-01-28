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
	"fmt"
	store2 "github.com/houyi-tracing/houyi/cmd/sm/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	api_v1.UnimplementedDynamicStrategyManagerServer
	api_v1.UnimplementedStrategyManagerServer
	api_v1.UnimplementedTraceGraphManagerServer
	api_v1.UnimplementedEvaluatorManagerServer

	logger *zap.Logger
	sst    sst.SamplingStrategyTree
	tg     tg.TraceGraph
	store  store2.OperationStore
	eval   evaluator.Evaluator
}

func NewGrpcHandler(logger *zap.Logger,
	sst sst.SamplingStrategyTree,
	tg tg.TraceGraph,
	store store2.OperationStore,
	eval evaluator.Evaluator) *GrpcHandler {
	return &GrpcHandler{
		logger: logger,
		sst:    sst,
		tg:     tg,
		store:  store,
		eval:   eval,
	}
}

func (h *GrpcHandler) Promote(_ context.Context, request *api_v1.Operation) (*api_v1.NullRely, error) {
	h.logger.Debug("Received request to Promote", zap.String("request", request.String()))

	reply := &api_v1.NullRely{}
	if h.tg.IsEntry(request) {
		err := h.sst.Promote(request)
		return reply, err
	} else {
		return reply, fmt.Errorf("can not promote operation which is not entry operation")
	}
}

func (h *GrpcHandler) GetStrategy(_ context.Context, request *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	h.logger.Debug("Received request to GetStrategy", zap.String("request", request.String()))

	strategies := make([]*api_v1.DynamicPerOperationSampling, 0)
	serviceName := request.GetService()
	for _, op := range request.GetOperations() {
		operation := &api_v1.Operation{
			Service:   serviceName,
			Operation: op.GetName(),
		}
		if h.tg.IsEntry(operation) {
			h.store.UpToDate(operation, true, op.Qps)
			if strategy, err := h.sst.Generate(operation); err == nil && strategy != nil {
				strategies = append(strategies,
					h.toStrategyResp(operation, strategy.SamplingRate*h.store.QpsWeight(operation)))
			} else {
				return nil, err
			}
		} else {
			h.store.UpToDate(operation, false, op.Qps)
			// Give operations that is not an entry with sampling rate 1.0 so that when it become entry, tracing
			// system would quickly collect enough traces to know that (this operation has become entry).
			strategies = append(strategies, h.toStrategyResp(operation, 1.0))
		}
	}

	return &api_v1.StrategyResponse{
		StrategyType: api_v1.StrategyResponse_DYNAMIC,
		Strategy: &api_v1.StrategyResponse_Dynamic{
			Dynamic: &api_v1.DynamicSampling{
				Strategies: strategies,
			},
		},
	}, nil
}

func (h *GrpcHandler) Traces(_ context.Context, operation *api_v1.Operation) (*api_v1.TracesReply, error) {
	h.logger.Debug("Received request to Traces", zap.String("operation", operation.String()))

	traces, err := h.tg.Traces(operation)
	return &api_v1.TracesReply{
		Entries: traces,
	}, err
}

func (h *GrpcHandler) GetTags(_ context.Context, _ *api_v1.GetTagsRequest) (*api_v1.GetTagsRely, error) {
	h.logger.Debug("Received request to GetTags")

	return &api_v1.GetTagsRely{Tags: h.eval.Get().Tags}, nil
}

func (h *GrpcHandler) UpdateTags(_ context.Context, request *api_v1.UpdateTagsRequest) (*api_v1.UpdateTagsReply, error) {
	h.logger.Debug("Received request to UpdateTags")

	h.eval.Update(&api_v1.EvaluatingTags{Tags: request.Tags})
	return &api_v1.UpdateTagsReply{}, nil
}

func (h *GrpcHandler) toStrategyResp(op *api_v1.Operation, samplingRate float64) *api_v1.DynamicPerOperationSampling {
	return &api_v1.DynamicPerOperationSampling{
		Operation: op.Operation,
		Strategy: &api_v1.DynamicPerOperationSampling_Probability{
			Probability: &api_v1.ProbabilitySampling{
				SamplingRate: samplingRate,
			},
		},
	}
}
