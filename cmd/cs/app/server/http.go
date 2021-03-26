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

package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	handler "github.com/houyi-tracing/houyi/cmd/cs/app/handler/http"
	"github.com/houyi-tracing/houyi/cmd/cs/app/store"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
)

type HttpServerParams struct {
	ListenPort int

	Logger *zap.Logger

	TraceGraph tg.TraceGraph

	StrategyStore store.StrategyStore

	Evaluator evaluator.Evaluator

	GossipRegistry gossip.Registry
}

func StartHttpServer(params *HttpServerParams) error {
	c := gin.Default()

	tHandler := handler.NewTraceGraphHttpHandler(&handler.TraceGraphHttpHandlerParams{
		Logger:     params.Logger,
		TraceGraph: params.TraceGraph,
	})
	tHandler.RegisterRoutes(c)

	smHandler := handler.NewStrategyManagerHttpHandler(&handler.StrategyManagerHttpHandlerParams{
		Logger:        params.Logger,
		StrategyStore: params.StrategyStore,
	})
	smHandler.RegisterRoutes(c)

	eHandler := handler.NewEvaluatorHttpHandler(&handler.EvaluatorHttpHandlerParams{
		Logger:         params.Logger,
		Evaluator:      params.Evaluator,
		GossipRegistry: params.GossipRegistry,
	})
	eHandler.RegisterRoutes(c)

	go func() {
		err := c.Run(fmt.Sprintf(":%d", params.ListenPort))
		if err != nil {
			params.Logger.Fatal("failed to run HTTP server", zap.Error(err))
		}
	}()

	return nil
}
