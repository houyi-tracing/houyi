// Copyright (c) 2019 The Jaeger Authors.
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

package server

import (
	app "github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	handler2 "github.com/houyi-tracing/houyi/cmd/collector/app/handler"
	"github.com/houyi-tracing/houyi/cmd/collector/app/sampling"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/collector/app/handler"
	clientcfgHandler "github.com/jaegertracing/jaeger/pkg/clientcfg/clientcfghttp"
	"github.com/jaegertracing/jaeger/pkg/healthcheck"
	"github.com/jaegertracing/jaeger/pkg/recoveryhandler"
)

// HTTPServerParams to construct a new Jaeger Collector HTTP Server
type HTTPServerParams struct {
	HostPort             string
	JaegerBatchesHandler handler.JaegerBatchesHandler
	SpanFilter           app.SpanFilter
	SamplingStore        sampling.AdaptiveStrategyStore
	MetricsFactory       metrics.Factory
	HealthCheck          *healthcheck.HealthCheck
	Logger               *zap.Logger
}

// StartHTTPServer based on the given parameters
func StartHTTPServer(params *HTTPServerParams) (*http.Server, error) {
	params.Logger.Info("Starting jaeger-collector HTTP server", zap.String("http host-port", params.HostPort))

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: params.HostPort}
	serveHTTP(server, listener, params)

	return server, nil
}

func serveHTTP(server *http.Server, listener net.Listener, params *HTTPServerParams) {
	r := mux.NewRouter()
	apiHandler := handler2.NewAPIHandler(params.JaegerBatchesHandler, params.SamplingStore, params.SpanFilter)
	apiHandler.RegisterRoutes(r)

	cfgHandler := clientcfgHandler.NewHTTPHandler(clientcfgHandler.HTTPHandlerParams{
		ConfigManager: &clientcfgHandler.ConfigManager{
			SamplingStrategyStore: params.SamplingStore,
		},
		MetricsFactory:         params.MetricsFactory,
		BasePath:               "/api",
		LegacySamplingEndpoint: false,
	})
	cfgHandler.RegisterRoutes(r)

	recoveryHandler := recoveryhandler.NewRecoveryHandler(params.Logger, true)
	server.Handler = recoveryHandler(r)
	go func() {
		if err := server.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				params.Logger.Fatal("Could not start HTTP collector", zap.Error(err))
			}
		}
		params.HealthCheck.Set(healthcheck.Unavailable)
	}()
}
