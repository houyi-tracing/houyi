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

package evaluator

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestMustReturnsFalseWhenGetDifferentTypeButSameTagName(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	eval := NewEvaluator(logger)
	tagName := "tag name"
	span := &model.Span{
		Tags: []model.KeyValue{
			{
				Key:   tagName,
				VType: model.ValueType_BOOL,
				VBool: true,
			},
		},
	}

	evaluatingTags := &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_INTEGER,
				Value: &api_v1.EvaluatingTag_IntegerVal{
					IntegerVal: 0,
				},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.False(t, eval.Evaluate(span))

	evaluatingTags = &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_STRING,
				Value: &api_v1.EvaluatingTag_StringVal{
					StringVal: "",
				},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.False(t, eval.Evaluate(span))

	evaluatingTags = &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_FLOAT,
				Value: &api_v1.EvaluatingTag_FloatVal{
					FloatVal: 0,
				},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.False(t, eval.Evaluate(span))
}

func TestMustReturnsTrueWhenTypeAndValueMatched(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	eval := NewEvaluator(logger)
	tagName := "tag name"

	span := &model.Span{
		Tags: []model.KeyValue{
			{
				Key:   tagName,
				VType: model.ValueType_BOOL,
				VBool: true,
			},
		},
	}
	evaluatingTags := &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_BOOLEAN,
				Value: &api_v1.EvaluatingTag_BooleanVal{
					BooleanVal: true},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.True(t, eval.Evaluate(span))

	span = &model.Span{
		Tags: []model.KeyValue{
			{
				Key:    tagName,
				VType:  model.ValueType_INT64,
				VInt64: 200,
			},
		},
	}
	evaluatingTags = &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_INTEGER,
				Value: &api_v1.EvaluatingTag_IntegerVal{
					IntegerVal: 200},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.True(t, eval.Evaluate(span))

	span = &model.Span{
		Tags: []model.KeyValue{
			{
				Key:   tagName,
				VType: model.ValueType_STRING,
				VStr:  "hello",
			},
		},
	}
	evaluatingTags = &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_STRING,
				Value: &api_v1.EvaluatingTag_StringVal{
					StringVal: "hello"},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.True(t, eval.Evaluate(span))

	span = &model.Span{
		Tags: []model.KeyValue{
			{
				Key:      tagName,
				VType:    model.ValueType_FLOAT64,
				VFloat64: 0.5,
			},
		},
	}
	evaluatingTags = &api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       tagName,
				OperationType: api_v1.EvaluatingTag_EQUAL_TO,
				ValueType:     api_v1.EvaluatingTag_FLOAT,
				Value: &api_v1.EvaluatingTag_FloatVal{
					FloatVal: 0.5},
			},
		},
	}
	eval.Update(evaluatingTags)
	assert.True(t, eval.Evaluate(span))
}
