// Copyright (c) 2020 The Houyi Authors.
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

package filter

import (
	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestSpanFilter(t *testing.T) {
	configFile := "filter-config.json"
	tsf := newTagSpanFilterFromFile(configFile, zap.NewNop())

	span := &model.Span{}
	assert.False(t, tsf.Filter(span))

	tags := []model.KeyValue{
		{
			Key:    "error",
			VInt64: 1,
			VType:  model.ValueType_INT64,
		},
	}
	span.Tags = tags
	assert.True(t, tsf.Filter(span))

	tags = []model.KeyValue{
		{
			Key:   "http.status_code",
			VStr:  "200",
			VType: model.ValueType_STRING,
		},
	}
	span.Tags = tags
	assert.False(t, tsf.Filter(span))

	tags = []model.KeyValue{
		{
			Key:    "request_size",
			VInt64: 300,
			VType:  model.ValueType_INT64,
		},
	}
	span.Tags = tags
	assert.False(t, tsf.Filter(span))

	tags = []model.KeyValue{
		{
			Key:    "body_length",
			VInt64: 500,
			VType:  model.ValueType_INT64,
		},
	}
	span.Tags = tags
	assert.True(t, tsf.Filter(span))
}

func TestJsonToFilterTags(t *testing.T) {
	jsonStr := `
	{
	  "filter-tags": [
		{
		  "key": "http.status_code",
		  "operation": "==",
		  "value": 200
		}
	  ]
	}
	`
	ft := JsonToFilterTags([]byte(jsonStr))
	assert.Equal(t, "http.status_code", ft.Tags[0].Key)
	assert.Equal(t, "==", ft.Tags[0].Operation)
	assert.Equal(t, int64(200), ft.Tags[0].Value)
}

func TestDefaultSpanFilter(t *testing.T) {
	logger, _ := zap.NewProduction()
	ft := newTagSpanFilterFromFile(defaultConfigFile, logger)
	assert.NotNil(t, ft)
}
