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
	opStore "github.com/houyi-tracing/houyi/cmd/sm/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	api_v1.UnimplementedDynamicStrategyManagerServer
	api_v1.UnimplementedStrategyManagerServer
	api_v1.UnimplementedTraceGraphManagerServer
	api_v1.UnimplementedEvaluatorManagerServer

	logger      *zap.Logger
	scaleFactor *atomic.Float64
	sst         sst.SamplingStrategyTree
	tg          tg.TraceGraph
	store       opStore.OperationStore
	eval        evaluator.Evaluator
	seed        gossip.Seed
}

func NewGrpcHandler(logger *zap.Logger,
	sst sst.SamplingStrategyTree,
	tg tg.TraceGraph,
	store opStore.OperationStore,
	eval evaluator.Evaluator,
	seed gossip.Seed,
	scaleFactor *atomic.Float64) *GrpcHandler {
	return &GrpcHandler{
		logger:      logger,
		sst:         sst,
		tg:          tg,
		store:       store,
		eval:        eval,
		seed:        seed,
		scaleFactor: scaleFactor,
	}
}

func (h *GrpcHandler) Promote(_ context.Context, request *api_v1.Operation) (*api_v1.NullRely, error) {
	h.logger.Debug("Received request to Promote", zap.String("request", request.String()))

	reply := &api_v1.NullRely{}
	if h.tg.IsIngress(request) {
		err := h.sst.Promote(request)
		return reply, err
	} else {
		return reply, fmt.Errorf("can not promote operation which is not entry operation")
	}
}

func (h *GrpcHandler) GetStrategy(_ context.Context, request *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	h.logger.Debug("Received request to GetStrategy", zap.String("request", request.String()))

	switch request.GetStrategyType() {
	case api_v1.StrategyType_ADAPTIVE:
		return h.adaptiveStrategy(request)
	case api_v1.StrategyType_DYNAMIC:
		return h.dynamicStrategy(request)
	default:
		return h.defaultStrategy(), nil
	}
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

func (h *GrpcHandler) GetServices(_ context.Context, req *api_v1.GetServicesRequest) (*api_v1.GetServicesReply, error) {
	h.logger.Debug("Received request to GetServices")

	return &api_v1.GetServicesReply{Services: h.tg.Services()}, nil
}

func (h *GrpcHandler) GetOperations(_ context.Context, req *api_v1.GetOperationsRequest) (*api_v1.GetOperationsReply, error) {
	h.logger.Debug("Received request to GetServices")

	return &api_v1.GetOperationsReply{Operations: h.tg.Operations(req.GetService())}, nil
}

func (h *GrpcHandler) dynamicStrategy(request *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	strategies := make([]*api_v1.DynamicPerOperationSampling, 0)
	serviceName := request.GetService()
	for _, op := range request.GetOperations() {
		operation := &api_v1.Operation{
			Service:   serviceName,
			Operation: op.GetName(),
		}
		if !h.tg.Has(operation) {
			_ = h.tg.Add(operation)
			h.seed.MongerNewOperation(operation)
		}
		if h.tg.IsIngress(operation) {
			h.store.UpToDate(operation, true, op.Qps)
			if strategy, err := h.sst.Generate(operation); err == nil && strategy != nil {
				strategies = append(strategies,
					h.toDynamicStrategyResp(operation, strategy.SamplingRate*h.store.QpsWeight(operation)))
			} else {
				return nil, err
			}
		} else {
			h.store.UpToDate(operation, false, op.Qps)
			// Give operations that is not an ingress with sampling rate 1.0 so that when it become entry, tracing
			// system would quickly collect enough traces to know that (this operation has become ingress).
			strategies = append(strategies, h.toDynamicStrategyResp(operation, 1.0))
		}
	}

	return &api_v1.StrategyResponse{
		StrategyType: api_v1.StrategyType_DYNAMIC,
		Strategy: &api_v1.StrategyResponse_Dynamic{
			Dynamic: &api_v1.DynamicSampling{
				Strategies: strategies,
			},
		},
	}, nil
}

func (h *GrpcHandler) adaptiveStrategy(request *api_v1.StrategyRequest) (*api_v1.StrategyResponse, error) {
	strategies := make([]*api_v1.PerOperationSampling, 0)
	serviceName := request.GetService()

	for _, op := range request.GetOperations() {
		operation := &api_v1.Operation{
			Service:   serviceName,
			Operation: op.GetName(),
		}
		if !h.tg.Has(operation) {
			_ = h.tg.Add(operation)
			h.seed.MongerNewOperation(operation)
		}
		if h.tg.IsIngress(operation) {
			h.store.UpToDate(operation, true, op.Qps)
			strategies = append(strategies,
				h.toAdaptiveStrategyResp(operation, h.store.QpsWeight(operation)))
		} else {
			h.store.UpToDate(operation, false, op.Qps)
			// Give operations that is not an ingress with sampling rate 1.0 so that when it become ingress, tracing
			// system would quickly collect enough traces to know that (this operation has become ingress).
			strategies = append(strategies, h.toAdaptiveStrategyResp(operation, 1.0))
		}
	}

	return &api_v1.StrategyResponse{
		StrategyType: api_v1.StrategyType_ADAPTIVE,
		Strategy: &api_v1.StrategyResponse_Adaptive{
			Adaptive: &api_v1.AdaptiveSampling{
				Strategies: strategies,
			},
		},
	}, nil
}

func (h *GrpcHandler) toDynamicStrategyResp(op *api_v1.Operation, samplingRate float64) *api_v1.DynamicPerOperationSampling {
	return &api_v1.DynamicPerOperationSampling{
		Operation: op.Operation,
		Strategy: &api_v1.DynamicPerOperationSampling_Probability{
			Probability: &api_v1.ProbabilitySampling{SamplingRate: samplingRate},
		},
	}
}

func (h *GrpcHandler) toAdaptiveStrategyResp(op *api_v1.Operation, samplingRate float64) *api_v1.PerOperationSampling {
	return &api_v1.PerOperationSampling{
		Operation: op.GetOperation(),
		Strategy:  &api_v1.ProbabilitySampling{SamplingRate: samplingRate},
	}
}

func (h *GrpcHandler) defaultStrategy() *api_v1.StrategyResponse {
	return &api_v1.StrategyResponse{
		StrategyType: api_v1.StrategyType_PROBABILITY,
		Strategy: &api_v1.StrategyResponse_Probability{
			Probability: &api_v1.ProbabilitySampling{SamplingRate: 1.0},
		},
	}
}
