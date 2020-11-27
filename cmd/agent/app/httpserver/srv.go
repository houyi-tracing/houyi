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

package httpserver

import (
	handler2 "github.com/yan-fuhai/adaptive-sampling-tracing/cmd/agent/app/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/jaegertracing/jaeger/cmd/agent/app/configmanager"
	"github.com/jaegertracing/jaeger/pkg/clientcfg/clientcfghttp"
)

// NewHTTPServer creates a new server that hosts an HTTP/JSON endpoint for clients
// to query for sampling strategies and baggage restrictions.
func NewHTTPServer(hostPort string, collectorHTTPHostPort string, manager configmanager.ClientConfigManager, mFactory metrics.Factory) *http.Server {
	handler := clientcfghttp.NewHTTPHandler(clientcfghttp.HTTPHandlerParams{
		ConfigManager:          manager,
		MetricsFactory:         mFactory,
		LegacySamplingEndpoint: true,
	})

	samplingHandler := handler2.NewSamplingHttpHandler(handler2.SamplingHttpHandlerParams{
		CollectorHostPort: collectorHTTPHostPort,
		MetricsFactory:    mFactory,
	})

	r := mux.NewRouter()

	handler.RegisterRoutes(r)
	samplingHandler.RegisterRoutes(r)

	return &http.Server{Addr: hostPort, Handler: r}
}
