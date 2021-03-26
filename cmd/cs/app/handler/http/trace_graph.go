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
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/houyi-tracing/houyi/route"
	"go.uber.org/zap"
	"net/http"
)

type TraceGraphHttpHandlerParams struct {
	Logger     *zap.Logger
	TraceGraph tg.TraceGraph
}

type TraceGraphHttpHandler struct {
	logger *zap.Logger
	tg     tg.TraceGraph
}

func NewTraceGraphHttpHandler(params *TraceGraphHttpHandlerParams) *TraceGraphHttpHandler {
	return &TraceGraphHttpHandler{
		logger: params.Logger,
		tg:     params.TraceGraph,
	}
}

func (h *TraceGraphHttpHandler) RegisterRoutes(e *gin.Engine) {
	e.GET(route.GetServicesRoute, h.getServices)
	e.GET(route.GetOperationsRoute, h.getOperations)
	e.GET(route.GetCausalDependenciesRoute, h.getCausalDependencies)
}

func (h *TraceGraphHttpHandler) getServices(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, gin.H{
		"result": h.tg.Services(),
	})
}

func (h *TraceGraphHttpHandler) getOperations(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	svc := c.Query("service")
	if svc != "" {
		h.logger.Debug("getOperations", zap.String("service name", svc))

		c.JSON(http.StatusOK, gin.H{
			"result": h.tg.Operations(svc),
		})
	} else {
		h.logger.Debug("getOperations: service is empty")

		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func (h *TraceGraphHttpHandler) getCausalDependencies(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	svc := c.Query("service")
	op := c.Query("operation")
	if svc == "" {
		h.logger.Debug("getEvaluatorTags: Empty service")
		c.Status(http.StatusBadRequest)
		return
	}
	if op == "" {
		h.logger.Debug("getEvaluatorTags: Empty operation")
		c.Status(http.StatusBadRequest)
		return
	}

	dependencies, err := h.tg.Dependencies(&api_v1.Operation{
		Service:   svc,
		Operation: op,
	})
	if err != nil {
		h.logger.Error("failed to get dependencies", zap.Error(err))
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": dependencies,
		})
	}
}

func (h *TraceGraphHttpHandler) getIngressServices(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, gin.H{
		"result": h.tg.AllIngresses(),
	})
}
