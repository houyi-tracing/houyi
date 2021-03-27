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
	c.GET(route.GetDefaultStrategyRoute, h.getDefaultStrategy)

	c.POST(route.UpdateStrategiesRoute, h.updateStrategies)
	c.POST(route.UpdateDefaultStrategyRoute, h.updateDefaultStrategy)
}

func (h *StrategyManagerHttpHandler) getStrategies(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	toResp := make([]*model.Strategy, 0)
	for _, s := range h.store.GetAll() {
		toResp = append(toResp, convertStrategyToJsonModel(s))
	}
	c.JSON(http.StatusOK, gin.H{
		"result": toResp,
	})
}

func (h *StrategyManagerHttpHandler) updateStrategies(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	strategies := make([]*model.Strategy, 0)
	if err := c.BindJSON(&strategies); err == nil {
		toUpdate := make([]*api_v1.PerOperationStrategy, 0)
		for _, s := range strategies {
			toUpdate = append(toUpdate, convertJsonModelToStrategy(s))
		}
		h.store.UpdateAll(toUpdate)
		c.JSON(http.StatusOK, gin.H{
			"result": "OK",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
		})
	}
}

func (h *StrategyManagerHttpHandler) getDefaultStrategy(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, gin.H{
		"result": convertStrategyToJsonModel(h.store.GetDefaultStrategy()),
	})
}

func (h *StrategyManagerHttpHandler) updateDefaultStrategy(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	newOne := &model.Strategy{}
	if err := c.BindJSON(newOne); err == nil {
		h.store.SetDefaultStrategy(convertJsonModelToStrategy(newOne))
		c.JSON(http.StatusOK, gin.H{
			"result": "OK",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
		})
	}
}

func convertStrategyToJsonModel(s *api_v1.PerOperationStrategy) *model.Strategy {
	ret := &model.Strategy{
		Service:   s.GetService(),
		Operation: s.GetOperation(),
	}

	switch s.GetType() {
	case api_v1.Type_CONST:
		ret.Type = model.StrategyType_Const
		ret.AlwaysSample = s.GetConst().GetAlwaysSample()
	case api_v1.Type_PROBABILITY:
		ret.Type = model.StrategyType_Probability
		ret.SamplingRate = s.GetProbability().GetSamplingRate()
	case api_v1.Type_RATE_LIMITING:
		ret.Type = model.StrategyType_RateLimiting
		ret.MaxTracesPerSecond = s.GetRateLimiting().GetMaxTracesPerSecond()
	case api_v1.Type_ADAPTIVE:
		ret.Type = model.StrategyType_Adaptive
	case api_v1.Type_DYNAMIC:
		ret.Type = model.StrategyType_Dynamic
	}

	return ret
}

func convertJsonModelToStrategy(s *model.Strategy) *api_v1.PerOperationStrategy {
	ret := &api_v1.PerOperationStrategy{
		Service:   s.Service,
		Operation: s.Operation,
	}

	switch s.Type {
	case model.StrategyType_Const:
		ret.Type = api_v1.Type_CONST
		ret.Strategy = &api_v1.PerOperationStrategy_Const{
			Const: &api_v1.ConstSampling{
				AlwaysSample: s.AlwaysSample,
			}}
	case model.StrategyType_Probability:
		ret.Type = api_v1.Type_PROBABILITY
		ret.Strategy = &api_v1.PerOperationStrategy_Probability{
			Probability: &api_v1.ProbabilitySampling{
				SamplingRate: s.SamplingRate,
			}}
	case model.StrategyType_RateLimiting:
		ret.Type = api_v1.Type_RATE_LIMITING
		ret.Strategy = &api_v1.PerOperationStrategy_RateLimiting{
			RateLimiting: &api_v1.RateLimitingSampling{
				MaxTracesPerSecond: s.MaxTracesPerSecond,
			}}
	case model.StrategyType_Adaptive:
		ret.Type = api_v1.Type_ADAPTIVE
	case model.StrategyType_Dynamic:
		ret.Type = api_v1.Type_DYNAMIC
	}

	return ret
}
