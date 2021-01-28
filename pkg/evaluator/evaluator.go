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
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"log"
	"sync"
)

// spanEvaluator filters spans by tags
type spanEvaluator struct {
	sync.RWMutex

	logger *zap.Logger
	tags   *api_v1.EvaluatingTags
	eqTags map[string]interface{} // equal to
	neTags map[string]interface{} // not equal to
	ltTags map[string]interface{} // less than
	gtTags map[string]interface{} // greater than
	leTags map[string]interface{} // less than or equal to
	geTags map[string]interface{} // greater than or equal to
}

func NewEvaluator(logger *zap.Logger) Evaluator {
	return &spanEvaluator{
		logger: logger,
		tags: &api_v1.EvaluatingTags{
			Tags: []*api_v1.EvaluatingTag{},
		},
		eqTags: make(map[string]interface{}),
		neTags: make(map[string]interface{}),
		ltTags: make(map[string]interface{}),
		gtTags: make(map[string]interface{}),
		leTags: make(map[string]interface{}),
		geTags: make(map[string]interface{}),
	}
}

func (f *spanEvaluator) Evaluate(span *model.Span) bool {
	f.RLock()
	defer f.RUnlock()

	for _, t := range span.GetTags() {
		switch t.GetVType() {
		case model.ValueType_BOOL:
			if f.checkBool(t.GetKey(), t.GetVBool()) {
				return true
			}
		case model.ValueType_FLOAT64:
			if f.checkFloat64(t.GetKey(), t.GetVFloat64()) {
				return true
			}
		case model.ValueType_STRING:
			if f.checkString(t.GetKey(), t.GetVStr()) {
				return true
			}
		case model.ValueType_INT64:
			if f.checkInt64(t.GetKey(), t.GetVInt64()) {
				return true
			}
		default:
			log.Println("unsupported Tag type:", t.GetVType())
		}
	}
	return false
}

func (f *spanEvaluator) Update(tags *api_v1.EvaluatingTags) {
	f.Lock()
	defer f.Unlock()

	f.clear()
	f.tags = tags
	f.parseTags(tags)
}

func (f *spanEvaluator) Get() *api_v1.EvaluatingTags {
	f.RLock()
	defer f.RUnlock()

	return f.tags
}

func (f *spanEvaluator) parseTags(tags *api_v1.EvaluatingTags) {
	for _, tag := range tags.Tags {
		switch tag.OperationType {
		case api_v1.EvaluatingTag_EQUAL_TO:
			f.eqTags[tag.TagName] = toActualType(tag)
		case api_v1.EvaluatingTag_NOT_EQUAL_TO:
			f.neTags[tag.TagName] = toActualType(tag)
		case api_v1.EvaluatingTag_GREATER_THAN:
			f.gtTags[tag.TagName] = toActualType(tag)
		case api_v1.EvaluatingTag_GREATER_THAN_OR_EQUAL_TO:
			f.geTags[tag.TagName] = toActualType(tag)
		case api_v1.EvaluatingTag_LESS_THAN:
			f.ltTags[tag.TagName] = toActualType(tag)
		case api_v1.EvaluatingTag_LESS_THAN_OR_EQUAL_TO:
			f.leTags[tag.TagName] = toActualType(tag)
		}
	}
}

// clear removes all tags.
func (f *spanEvaluator) clear() {
	f.tags = &api_v1.EvaluatingTags{Tags: []*api_v1.EvaluatingTag{}}
	f.eqTags = make(map[string]interface{})
	f.neTags = make(map[string]interface{})
	f.gtTags = make(map[string]interface{})
	f.geTags = make(map[string]interface{})
	f.ltTags = make(map[string]interface{})
	f.leTags = make(map[string]interface{})
}

func (f *spanEvaluator) checkBool(tKey string, tVal bool) bool {
	if cmp, has := f.eqTags[tKey]; has {
		if cVal, err := cast.ToBoolE(cmp); err == nil && cVal == tVal {
			return true
		}
	}
	if cmp, has := f.neTags[tKey]; has {
		if cVal, err := cast.ToBoolE(cmp); err == nil && cVal != tVal {
			return true
		}
	}
	return false
}

func (f *spanEvaluator) checkFloat64(tKey string, tVal float64) bool {
	if cmp, has := f.eqTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal == tVal {
			return true
		}
	}
	if cmp, has := f.neTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal != tVal {
			return true
		}
	}
	if cmp, has := f.ltTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal < tVal {
			return true
		}
	}
	if cmp, has := f.gtTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal > tVal {
			return true
		}
	}
	if cmp, has := f.leTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal <= tVal {
			return true
		}
	}
	if cmp, has := f.geTags[tKey]; has {
		if cVal, err := cast.ToFloat64E(cmp); err == nil && cVal >= tVal {
			return true
		}
	}
	return false
}

func (f *spanEvaluator) checkString(tKey string, tVal string) bool {
	if cmp, has := f.eqTags[tKey]; has {
		if cVal, err := cast.ToStringE(cmp); err == nil && cVal == tVal {
			return true
		}
	}
	if cmp, has := f.neTags[tKey]; has {
		if cVal, err := cast.ToStringE(cmp); err == nil && cVal != tVal {
			return true
		}
	}
	return false
}

func (f *spanEvaluator) checkInt64(tKey string, tVal int64) bool {
	if cmp, has := f.eqTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal == tVal {
			return true
		}
	}
	if cmp, has := f.neTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal != tVal {
			return true
		}
	}
	if cmp, has := f.ltTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal < tVal {
			return true
		}
	}
	if cmp, has := f.gtTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal > tVal {
			return true
		}
	}
	if cmp, has := f.leTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal <= tVal {
			return true
		}
	}
	if cmp, has := f.geTags[tKey]; has {
		if cVal, err := cast.ToInt64E(cmp); err == nil && cVal >= tVal {
			return true
		}
	}
	return false
}

func toActualType(tag *api_v1.EvaluatingTag) interface{} {
	switch tag.ValueType {
	case api_v1.EvaluatingTag_BOOLEAN:
		return tag.GetBooleanVal()
	case api_v1.EvaluatingTag_INTEGER:
		return tag.GetIntegerVal()
	case api_v1.EvaluatingTag_FLOAT:
		return tag.GetFloatVal()
	case api_v1.EvaluatingTag_STRING:
		return tag.GetStringVal()
	}
	return nil
}
