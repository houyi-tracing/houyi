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

package processor

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test(t *testing.T) {
	spans := []*model.Span{
		{
			TraceID:       model.TraceID{},
			SpanID:        0,
			OperationName: "",
			References:    nil,
			Flags:         0,
			StartTime:     time.Time{},
			Duration:      0,
			Tags:          nil,
			Logs:          nil,
			Process:       nil,
			ProcessID:     "",
			Warnings:      nil,
		},
	}

	logger, _ := zap.NewDevelopment()
	sp := NewSpanProcessor(logger)

	err := sp.ProcessSpans(spans)
	assert.Nil(t, err)
}
