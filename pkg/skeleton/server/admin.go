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

package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/houyi-tracing/houyi/pkg/hc"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net"
	"net/http"
)

const (
	adminHttpPort = "admin.http.port"
	hcRoute       = "hc.route"

	DefaultHcRoute = "/health"
)

type AdminServer struct {
	httpPort int

	hcRoute string

	hc hc.HealthCheck

	logger *zap.Logger

	mux *http.ServeMux

	server *http.Server
}

func NewAdminServer(httpPort int) *AdminServer {
	return &AdminServer{
		httpPort: httpPort,
		hcRoute:  DefaultHcRoute,
		hc:       hc.NewHC(zap.NewNop()),
		logger:   zap.NewNop(),
		mux:      http.NewServeMux(),
	}
}

func (s *AdminServer) AddFlags(flags *flag.FlagSet) {
	flags.Int(adminHttpPort, ports.AdminHttpPort, "Port to serve Admin server")
	flags.String(hcRoute, DefaultHcRoute, "Route for serving health check")
}

func (s *AdminServer) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.httpPort))
	if err != nil {
		s.logger.Error("Admin server failed to listen", zap.Error(err))
		return err
	}
	s.serveHttp(l)
	return nil
}

func (s *AdminServer) HC() hc.HealthCheck {
	return s.hc
}

func (s *AdminServer) Close() error {
	return s.server.Shutdown(context.Background())
}

func (s *AdminServer) setLogger(logger *zap.Logger) {
	s.logger = logger
	s.hc.SetLogger(logger)
}

func (s *AdminServer) InitFromViper(v *viper.Viper, logger *zap.Logger) {
	s.setLogger(logger)
	s.httpPort = v.GetInt(adminHttpPort)
	s.hcRoute = v.GetString(hcRoute)
}

func (s *AdminServer) serveHttp(l net.Listener) {
	s.logger.Info("Starting health check on admin server", zap.String("route", s.hcRoute), zap.Int("http-port", s.httpPort))

	s.mux.Handle(s.hcRoute, s.hc.Handler())
	s.server = &http.Server{Handler: http.Handler(s.mux)}
	go func() {
		if err := s.server.Serve(l); err != nil {
			s.logger.Info("Admin server stopped", zap.Error(err))
			s.hc.Set(hc.Broken)
		}
	}()
}
