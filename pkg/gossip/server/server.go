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
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip"
	"github.com/houyi-tracing/houyi/pkg/gossip/handler"
	"github.com/houyi-tracing/houyi/pkg/gossip/seed"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
)

type SeedParams struct {
	Logger               *zap.Logger
	ListenPort           int
	LruSize              int
	ConfigServerEndpoint *routing.Endpoint
	TraceGraph           tg.TraceGraph
	Evaluator            evaluator.Evaluator
}

func BuildSeed(params *SeedParams) (gossip.Seed, error) {
	s := seed.NewSeed(
		params.Logger,
		seed.Options.ListenPort(params.ListenPort),
		seed.Options.LruSize(params.LruSize),
		seed.Options.ConfigServerEndpoint(params.ConfigServerEndpoint))

	gHandler := handler.NewHandler(params.Logger, params.TraceGraph, params.Evaluator)

	s.OnNewOperation(gHandler.NewOperationHandler)
	s.OnExpiredOperation(gHandler.ExpiredOperationHandler)
	s.OnEvaluatingTags(gHandler.TagsHandler)
	s.OnNewRelation(gHandler.RelationHandler)

	return s, nil
}
