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

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/houyi-tracing/houyi/cmd/cs/app/handler/http/model"
	"github.com/houyi-tracing/houyi/cmd/cs/app/store"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/route"
	"go.uber.org/zap"
	"net/http"
)

type StrategyManagerHttpHandlerParams struct {
	Logger        *zap.Logger
	StrategyStore store.StrategyStore
}

type StrategyManagerHttpHandler struct {
	logger *zap.Logger
	store  store.StrategyStore
}

func NewStrategyManagerHttpHandler(params *StrategyManagerHttpHandlerParams) *StrategyManagerHttpHandler {
	return &StrategyManagerHttpHandler{
		logger: params.Logger,
		store:  params.StrategyStore,
	}
}

func (h *StrategyManagerHttpHandler) RegisterRoutes(c *gin.Engine) {
	c.GET(route.GetStrategiesRoute, h.getStrategies)
	c.POST(route.UpdateStrategiesRoute, h.updateStrategies)
}

func (h *StrategyManagerHttpHandler) getStrategies(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, gin.H{
		"result": convertStrategyToJsonModel(h.store.GetAll()),
	})
}

func (h *StrategyManagerHttpHandler) updateStrategies(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	strategies := make([]*model.Strategy, 0)
	if err := c.BindJSON(&strategies); err == nil {
		h.store.UpdateAll(convertJsonModelToStrategy(strategies))
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func convertStrategyToJsonModel(strategies []*api_v1.PerOperationStrategy) []*model.Strategy {
	ret := make([]*model.Strategy, 0)

	for _, s := range strategies {
		newS := &model.Strategy{
			Service:   s.GetService(),
			Operation: s.GetOperation(),
		}

		switch s.GetType() {
		case api_v1.Type_CONST:
			newS.Type = model.StrategyType_Const
			newS.AlwaysSample = s.GetConst().GetAlwaysSample()
		case api_v1.Type_PROBABILITY:
			newS.Type = model.StrategyType_Probability
			newS.SamplingRate = s.GetProbability().GetSamplingRate()
		case api_v1.Type_RATE_LIMITING:
			newS.Type = model.StrategyType_RateLimiting
			newS.MaxTracesPerSecond = s.GetRateLimiting().GetMaxTracesPerSecond()
		case api_v1.Type_ADAPTIVE:
			newS.Type = model.StrategyType_Adaptive
		case api_v1.Type_DYNAMIC:
			newS.Type = model.StrategyType_Dynamic
		}

		ret = append(ret, newS)
	}

	return ret
}

func convertJsonModelToStrategy(strategies []*model.Strategy) []*api_v1.PerOperationStrategy {
	ret := make([]*api_v1.PerOperationStrategy, 0)

	for _, s := range strategies {
		newS := &api_v1.PerOperationStrategy{
			Service:   s.Service,
			Operation: s.Operation,
		}

		switch s.Type {
		case model.StrategyType_Const:
			newS.Type = api_v1.Type_CONST
			newS.Strategy = &api_v1.PerOperationStrategy_Const{
				Const: &api_v1.ConstSampling{
					AlwaysSample: s.AlwaysSample,
				}}
		case model.StrategyType_Probability:
			newS.Type = api_v1.Type_PROBABILITY
			newS.Strategy = &api_v1.PerOperationStrategy_Probability{
				Probability: &api_v1.ProbabilitySampling{
					SamplingRate: s.SamplingRate,
				}}
		case model.StrategyType_RateLimiting:
			newS.Type = api_v1.Type_RATE_LIMITING
			newS.Strategy = &api_v1.PerOperationStrategy_RateLimiting{
				RateLimiting: &api_v1.RateLimitingSampling{
					MaxTracesPerSecond: s.MaxTracesPerSecond,
				}}
		case model.StrategyType_Adaptive:
			newS.Type = api_v1.Type_ADAPTIVE
		case model.StrategyType_Dynamic:
			newS.Type = api_v1.Type_DYNAMIC
		}

		ret = append(ret, newS)
	}

	return ret
}
