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
	"github.com/gorilla/mux"
	"github.com/uber/jaeger-lib/metrics"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPHandlerParams contains parameters that must be passed to NewHTTPHandler.
type SamplingHttpHandlerParams struct {
	CollectorHostPort string
	MetricsFactory    metrics.Factory // required
}

type SamplingHttpHandler struct {
	collectorHostPort string
	metricsFactory    metrics.Factory
}

func NewSamplingHttpHandler(params SamplingHttpHandlerParams) *SamplingHttpHandler {
	return &SamplingHttpHandler{
		collectorHostPort: params.CollectorHostPort,
		metricsFactory:    params.MetricsFactory,
	}
}

func (shh *SamplingHttpHandler) RegisterRoutes(route *mux.Router) {
	route.HandleFunc("/api/strategy", shh.getSamplingStrategy).Methods(http.MethodGet)
}

func (shh *SamplingHttpHandler) getSamplingStrategy(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://" + shh.collectorHostPort + r.URL.RequestURI())
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}
