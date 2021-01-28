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
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	tg     tg.TraceGraph
	eval   evaluator.Evaluator
}

func NewHandler(logger *zap.Logger, tg tg.TraceGraph, eval evaluator.Evaluator) *Handler {
	return &Handler{
		logger: logger,
		tg:     tg,
		eval:   eval,
	}
}

func (h *Handler) RelationHandler(rel *api_v1.Relation) {
	h.logger.Debug("Handle new relation", zap.String("relation", rel.String()))

	from, to := rel.From, rel.To
	if !h.tg.Has(from) {
		if err := h.tg.Add(from); err != nil {
			h.logger.Error("failed to add operation for trace graph", zap.Error(err))
		}
	}
	if !h.tg.Has(to) {
		if err := h.tg.Add(to); err != nil {
			h.logger.Error("failed to add operation for trace graph", zap.Error(err))
		}
	}
	if err := h.tg.AddRelation(rel); err != nil {
		h.logger.Error("failed to add relation for trace graph", zap.Error(err))
	}
}

func (h *Handler) OperationHandler(op *api_v1.Operation) {
	h.logger.Debug("Handle expired operation", zap.String("operation", op.String()))

	if h.tg.Has(op) {
		if err := h.tg.Remove(op); err != nil {
			h.logger.Error("failed to remove expired operation for trace graph", zap.Error(err))
		}
	}
}

func (h *Handler) TagsHandler(tags *api_v1.EvaluatingTags) {
	h.logger.Debug("Handle new evaluating tags", zap.String("evaluating tags", tags.String()))

	h.eval.Update(tags)
}
