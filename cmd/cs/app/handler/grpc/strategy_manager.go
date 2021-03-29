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

package grpc

import (
	"context"
	"github.com/houyi-tracing/houyi/cmd/cs/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
	"math"
)

type StrategyManagerGrpcHandler struct {
	api_v1.UnimplementedStrategyManagerServer

	logger          *zap.Logger
	scaleFactor     float64
	sst             sst.SamplingStrategyTree
	tg              tg.TraceGraph
	operationStore  store.OperationStore
	strategyStore   store.StrategyStore
	eval            evaluator.Evaluator
	gossipSeed      gossip.Seed
	minSamplingRate float64
}

func NewStrategyManagerGrpcHandler(logger *zap.Logger,
	sst sst.SamplingStrategyTree,
	tg tg.TraceGraph,
	opStore store.OperationStore,
	eval evaluator.Evaluator,
	scaleFactor float64,
	strategyStore store.StrategyStore,
	minSamplingRate float64,
	seed gossip.Seed) *StrategyManagerGrpcHandler {
	return &StrategyManagerGrpcHandler{
		logger:          logger,
		sst:             sst,
		tg:              tg,
		operationStore:  opStore,
		eval:            eval,
		scaleFactor:     scaleFactor,
		minSamplingRate: minSamplingRate,
		strategyStore:   strategyStore,
		gossipSeed:      seed,
	}
}

func (h *StrategyManagerGrpcHandler) Promote(_ context.Context, request *api_v1.Operation) (*api_v1.NullRely, error) {
	h.logger.Debug("Received request to Promote", zap.String("request", request.String()))

	reply := &api_v1.NullRely{}
	if h.tg.IsIngress(request) {
		err := h.sst.Promote(request)
		return reply, err
	} else {
		if ingress, err := h.tg.GetIngresses(request); err != nil {
			return reply, err
		} else {
			for _, i := range ingress {
				h.logger.Debug("Promoted operation",
					zap.String("service", i.GetService()),
					zap.String("operation", i.GetOperation()))
				err = h.sst.Promote(i)
			}
			return reply, err
		}
	}
}

func (h *StrategyManagerGrpcHandler) GetStrategies(_ context.Context, request *api_v1.StrategyRequest) (*api_v1.StrategiesResponse, error) {
	h.logger.Debug("Received request to GetStrategies", zap.String("request", request.String()))

	resp := &api_v1.StrategiesResponse{Strategies: make([]*api_v1.PerOperationStrategy, 0)}

	svc := request.GetService()
	for _, op := range request.GetOperations() {
		opModel := &api_v1.Operation{
			Service:   svc,
			Operation: op.GetName(),
		}
		isIngress := h.tg.IsIngress(opModel)
		h.operationStore.UpToDate(opModel, isIngress, op.GetQps())
		resp.Strategies = append(resp.Strategies, h.perOperationStrategy(opModel, isIngress))
	}
	return resp, nil
}

func (h *StrategyManagerGrpcHandler) perOperationStrategy(opModel *api_v1.Operation, isIngress bool) *api_v1.PerOperationStrategy {
	if !h.tg.Has(opModel) {
		h.gossipSeed.MongerNewOperation(opModel)
		_ = h.tg.Add(opModel)
	}

	svc, op := opModel.GetService(), opModel.GetOperation()
	if isIngress {
		ret := &api_v1.PerOperationStrategy{}
		if !h.strategyStore.Has(svc, op) {
			ret = h.strategyStore.GetDefaultStrategy()
		} else {
			ret, _ = h.strategyStore.Get(svc, op)
		}

		if ret.GetType() == api_v1.Type_DYNAMIC {
			if !h.sst.Has(opModel) {
				_ = h.sst.Add(opModel)
			}
			sr, _ := h.sst.Generate(opModel)
			qpsWeight := h.operationStore.QpsWeight(opModel)
			ret.Strategy = &api_v1.PerOperationStrategy_Dynamic{
				Dynamic: &api_v1.DynamicSampling{
					SamplingRate: math.Min(math.Max(sr*qpsWeight*h.scaleFactor, h.minSamplingRate), 1.0),
				}}
			h.logger.Debug("Generated dynamic strategy",
				zap.String("service", svc),
				zap.String("operation", op),
				zap.Float64("SST", sr),
				zap.Float64("QPS weight", qpsWeight))
		} else if ret.GetType() == api_v1.Type_ADAPTIVE {
			qpsWeight := h.operationStore.QpsWeight(opModel)
			ret.Strategy = &api_v1.PerOperationStrategy_Adaptive{
				Adaptive: &api_v1.AdaptiveSampling{
					SamplingRate: math.Min(math.Max(qpsWeight*h.scaleFactor, h.minSamplingRate), 1.0),
				}}
			h.logger.Debug("Generated adaptive strategy",
				zap.String("service", svc),
				zap.String("operation", op),
				zap.Float64("QPS weight", qpsWeight))
		}
		return ret
	} else {
		if err := h.sst.Prune(opModel); err == nil {
			h.logger.Debug("Removed non-ingress operation from SST.",
				zap.String("service", svc),
				zap.String("operation", op))
		}

		ret := h.strategyStore.GetDefaultStrategy()
		ret.Service = svc
		ret.Operation = op
		return ret
	}
}
