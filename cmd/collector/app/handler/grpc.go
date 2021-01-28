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
	"github.com/houyi-tracing/houyi/cmd/collector/app/processor"
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	logger        *zap.Logger
	spanProcessor processor.SpanProcessor
}

func NewGrpcHandler(logger *zap.Logger, sp processor.SpanProcessor) *GrpcHandler {
	return &GrpcHandler{
		logger:        logger,
		spanProcessor: sp,
	}
}

func (g *GrpcHandler) PostSpans(_ context.Context, request *api_v2.PostSpansRequest) (*api_v2.PostSpansResponse, error) {
	reply := &api_v2.PostSpansResponse{}
	if g.spanProcessor != nil {
		if err := g.spanProcessor.ProcessSpans(request.GetBatch().Spans); err != nil {
			return reply, nil
		} else {
			return reply, err
		}
	} else {
		return reply, fmt.Errorf("span processor is nil")
	}
}
