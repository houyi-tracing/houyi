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
	"encoding/json"
	"github.com/jaegertracing/jaeger/model"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	OperationTypeEq = "==" // equal to
	OperationTypeNe = "!=" // not equal to
	OperationTypeLt = "<"  // less than
	OperationTypeGt = ">"  // greater than
	OperationTypeLe = "<=" // less than or equal to
	OperationTypeGe = ">=" // greater than or equal to
)

// nullSpanFilter would not filter any spans
type nullSpanFilter struct{}

func NewNullSpanFilter() *nullSpanFilter {
	return &nullSpanFilter{}
}

func (nsf *nullSpanFilter) Filter(_ *model.Span) bool {
	return false
}

func (nsf *nullSpanFilter) Update(_ interface{}) {}

func (nsf *nullSpanFilter) GetConditions() interface{} {
	return nil
}

func (nsf *nullSpanFilter) Clear() {}

type Tags struct {
	Tags []Tag `json:"filter-tags"`
}

type Tag struct {
	Key       string      `json:"key"`
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}

// tagSpanFilter filters spans by tags
type tagSpanFilter struct {
	logger *zap.Logger
	eqTags map[string]interface{} // equal to
	neTags map[string]interface{} // not equal to
	ltTags map[string]interface{} // less than
	gtTags map[string]interface{} // greater than
	leTags map[string]interface{} // less than or equal to
	geTags map[string]interface{} // greater than or equal to
}

func newTagSpanFilterFromFile(configFile string, logger *zap.Logger) SpanFilter {
	filter := newTagSpanFilter(logger)
	if configFile == "" {
		filter.logger.Error("name of filter configuration file is empty")
		return filter
	}

	if configFile == defaultConfigFile {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			logger.Error("cannot read default filter configuration file because HOME env is empty")
			return filter
		} else {
			filePath := strings.Split(defaultConfigFile, "/")[1]
			configFile = homeDir + "/" + filePath
		}
	}

	logger.Info("reading filter configuration file", zap.String("file", configFile))
	jsonBytes, err := readAll(configFile)
	if err != nil {
		filter.logger.Error("failed to read filter tags file", zap.Error(err))
		return filter
	}

	filterTags := JsonToFilterTags(jsonBytes)
	filter.updateWithFilterTags(filterTags)
	return filter
}

func newTagSpanFilter(logger *zap.Logger) *tagSpanFilter {
	return &tagSpanFilter{
		logger: logger,
		eqTags: make(map[string]interface{}),
		neTags: make(map[string]interface{}),
		ltTags: make(map[string]interface{}),
		gtTags: make(map[string]interface{}),
		leTags: make(map[string]interface{}),
		geTags: make(map[string]interface{}),
	}
}

func (f *tagSpanFilter) Filter(span *model.Span) bool {
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

// Update updates the tags with inputted jsonStr string and it would not remove keys already exist in tags
// but not in inputted jsonStr.
func (f *tagSpanFilter) Update(v interface{}) {
	if filterTags, ok := v.(Tags); ok {
		f.updateWithFilterTags(&filterTags)
	}
}

// GetFilter returns the current filter tags this span filter using.
func (f *tagSpanFilter) GetConditions() interface{} {
	filterTags := Tags{}
	ops := []string{
		OperationTypeEq, OperationTypeNe, OperationTypeGt,
		OperationTypeGe, OperationTypeLt, OperationTypeLe}
	tags := []map[string]interface{}{
		f.eqTags, f.neTags, f.gtTags,
		f.geTags, f.ltTags, f.leTags}

	for i := 0; i < len(ops); i++ {
		for k, v := range tags[i] {
			filterTags.Tags = append(filterTags.Tags, Tag{
				Key:       k,
				Operation: ops[i],
				Value:     v,
			})
		}
	}
	return filterTags
}

// Clear removes all tags.
func (f *tagSpanFilter) Clear() {
	f.eqTags = make(map[string]interface{})
	f.neTags = make(map[string]interface{})
	f.gtTags = make(map[string]interface{})
	f.geTags = make(map[string]interface{})
	f.ltTags = make(map[string]interface{})
	f.leTags = make(map[string]interface{})
}

func (f *tagSpanFilter) updateWithFilterTags(filterTags *Tags) {
	for _, tag := range filterTags.Tags {
		f.removeKey(tag.Key)
		switch tag.Operation {
		case OperationTypeEq:
			f.eqTags[tag.Key] = tag.Value
		case OperationTypeNe:
			f.neTags[tag.Key] = tag.Value
		case OperationTypeLt:
			f.ltTags[tag.Key] = tag.Value
		case OperationTypeGt:
			f.gtTags[tag.Key] = tag.Value
		case OperationTypeLe:
			f.leTags[tag.Key] = tag.Value
		case OperationTypeGe:
			f.geTags[tag.Key] = tag.Value
		}
	}
}

func (f *tagSpanFilter) checkBool(tKey string, tVal bool) bool {
	if cmp, has := f.eqTags[tKey]; has && cast.ToBool(cmp) == tVal {
		return true
	}
	if cmp, has := f.neTags[tKey]; has && cast.ToBool(cmp) != tVal {
		return true
	}
	return false
}

func (f *tagSpanFilter) checkFloat64(tKey string, tVal float64) bool {
	if cmp, has := f.eqTags[tKey]; has && cast.ToFloat64(cmp) == tVal {
		return true
	}
	if cmp, has := f.neTags[tKey]; has && cast.ToFloat64(cmp) != tVal {
		return true
	}
	if cmp, has := f.ltTags[tKey]; has && cast.ToFloat64(cmp) < tVal {
		return true
	}
	if cmp, has := f.gtTags[tKey]; has && cast.ToFloat64(cmp) > tVal {
		return true
	}
	if cmp, has := f.leTags[tKey]; has && cast.ToFloat64(cmp) <= tVal {
		return true
	}
	if cmp, has := f.geTags[tKey]; has && cast.ToFloat64(cmp) >= tVal {
		return true
	}
	return false
}

func (f *tagSpanFilter) checkString(tKey string, tVal string) bool {
	if cmp, has := f.eqTags[tKey]; has && cast.ToString(cmp) == tVal {
		return true
	}
	if cmp, has := f.neTags[tKey]; has && cast.ToString(cmp) != tVal {
		return true
	}
	return false
}

func (f *tagSpanFilter) checkInt64(tKey string, tVal int64) bool {
	if cmp, has := f.eqTags[tKey]; has && cast.ToInt64(cmp) == tVal {
		return true
	}
	if cmp, has := f.neTags[tKey]; has && cast.ToInt64(cmp) != tVal {
		return true
	}
	if cmp, has := f.ltTags[tKey]; has && cast.ToInt64(cmp) < tVal {
		return true
	}
	if cmp, has := f.gtTags[tKey]; has && cast.ToInt64(cmp) > tVal {
		return true
	}
	if cmp, has := f.leTags[tKey]; has && cast.ToInt64(cmp) <= tVal {
		return true
	}
	if cmp, has := f.geTags[tKey]; has && cast.ToInt64(cmp) >= tVal {
		return true
	}
	return false
}

func (f *tagSpanFilter) removeKey(tKey string) {
	delete(f.eqTags, tKey)
	delete(f.neTags, tKey)
	delete(f.ltTags, tKey)
	delete(f.gtTags, tKey)
	delete(f.geTags, tKey)
	delete(f.leTags, tKey)
}

func readAll(filename string) ([]byte, error) {
	if ret, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

func isInt(fn float64) bool {
	return float64(int64(fn)) == fn
}

func JsonToFilterTags(bytes []byte) *Tags {
	ret := new(Tags)
	if err := json.Unmarshal(bytes, ret); err != nil {
		log.Println(err)
	}
	// package json would unmarshall all types of numbers in float numbers, including integer numbers.
	// Thus, we should convert those float numbers which is integer numbers numerically to integer numbers.
	for i, t := range ret.Tags {
		if fn, ok := t.Value.(float64); ok {
			if isInt(fn) {
				ret.Tags[i].Value = int64(fn)
			}
		}
	}
	return ret
}

func FilterTagsToJson(tags Tags) []byte {
	if bytes, err := json.Marshal(tags); err != nil {
		log.Println(err)
		return nil
	} else {
		return bytes
	}
}
