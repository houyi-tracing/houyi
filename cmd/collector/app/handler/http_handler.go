// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
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

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gorilla/mux"
	app "github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/cmd/collector/app/sampling"
	model2 "github.com/houyi-tracing/houyi/cmd/collector/app/sampling/model"
	"github.com/jaegertracing/jaeger/cmd/collector/app/handler"
	"github.com/jaegertracing/jaeger/cmd/collector/app/processor"
	"github.com/jaegertracing/jaeger/pkg/healthcheck"
	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"io/ioutil"
	"mime"
	"net/http"
	"time"
)

const (
	// UnableToReadBodyErrFormat is an error message for invalid requests
	UnableToReadBodyErrFormat = "Unable to process request body: %v"
)

var (
	acceptedThriftFormats = map[string]struct{}{
		"application/x-thrift":                 {},
		"application/vnd.apache.thrift.binary": {},
	}
)

// APIHandler handles all HTTP calls to the collector
type APIHandler struct {
	jaegerBatchesHandler handler.JaegerBatchesHandler
	strategyStore        sampling.AdaptiveStrategyStore
	spanFilter           app.SpanFilter
	hc                   *healthcheck.HealthCheck
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(
	jaegerBatchesHandler handler.JaegerBatchesHandler,
	strategyStore sampling.AdaptiveStrategyStore,
	filter app.SpanFilter,
	hc *healthcheck.HealthCheck,
) *APIHandler {
	return &APIHandler{
		jaegerBatchesHandler: jaegerBatchesHandler,
		strategyStore:        strategyStore,
		spanFilter:           filter,
		hc:                   hc,
	}
}

// RegisterRoutes registers routes for this handler on the given router
func (aH *APIHandler) RegisterRoutes(router *mux.Router) {
	// GET
	router.HandleFunc("/api/filter_tags", aH.GetFilterTags).Methods(http.MethodGet)
	router.HandleFunc("/api/strategy", aH.GetSamplingStrategy).Methods(http.MethodGet)
	router.Handle("/health", aH.hc.Handler())

	// POST
	router.HandleFunc("/api/traces", aH.SaveSpan).Methods(http.MethodPost)
	router.HandleFunc("/api/update_filter", aH.UpdateFilter).Methods(http.MethodPost)
}

// SaveSpan submits the span provided in the request body to the JaegerBatchesHandler
func (aH *APIHandler) SaveSpan(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf(UnableToReadBodyErrFormat, err), http.StatusInternalServerError)
		return
	}

	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse content type: %v", err), http.StatusBadRequest)
		return
	}

	if _, ok := acceptedThriftFormats[contentType]; !ok {
		http.Error(w, fmt.Sprintf("Unsupported content type: %v", contentType), http.StatusBadRequest)
		return
	}

	tdes := thrift.NewTDeserializer()
	batch := &tJaeger.Batch{}
	if err = tdes.Read(batch, bodyBytes); err != nil {
		http.Error(w, fmt.Sprintf(UnableToReadBodyErrFormat, err), http.StatusBadRequest)
		return
	}
	batches := []*tJaeger.Batch{batch}
	opts := handler.SubmitBatchOptions{InboundTransport: processor.HTTPTransport}
	if _, err = aH.jaegerBatchesHandler.SubmitBatches(batches, opts); err != nil {
		http.Error(w, fmt.Sprintf("Cannot submit Jaeger batch: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetSamplingStrategy process requests for get sampling strategy of single server.
// The QPS, which refers to "reQuest Per Second" must exist in the parameters.
func (aH *APIHandler) GetSamplingStrategy(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	operations := model2.Operations{}
	opParams := params["operation"]
	for _, sp := range opParams {
		op := model2.Operation{}
		if err := json.Unmarshal([]byte(sp), &op); err == nil {
			operations.Operations = append(operations.Operations, op)
		}
	}

	intervalStr := params.Get("interval")
	if intervalStr == "" {
		http.Error(w, "interval must be non-empty", http.StatusBadRequest)
		return
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serviceName := params.Get("service")
	if serviceName == "" {
		http.Error(w, "service name must be non-empty", http.StatusBadRequest)
		return
	}

	resp, err := aH.strategyStore.GetSamplingStrategies(serviceName, operations, interval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respToWrite, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respToWrite)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateFilter updates the filter tags for filtering spans.
func (aH *APIHandler) UpdateFilter(w http.ResponseWriter, r *http.Request) {
	newFilterTags, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if filterTags, err := parseFilterTags(newFilterTags); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		aH.spanFilter.Update(filterTags)
	}
}

func (aH *APIHandler) GetFilterTags(w http.ResponseWriter, _ *http.Request) {
	if resp := app.FilterTagsToJson(aH.spanFilter.GetConditions().(app.Tags)); resp != nil {
		_, _ = w.Write(resp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func parseFilterTags(data []byte) (app.Tags, error) {
	filterTags := app.Tags{}
	if err := json.Unmarshal(data, &filterTags); err != nil {
		return filterTags, fmt.Errorf("failed to parse filter tags")
	}
	return filterTags, nil
}
