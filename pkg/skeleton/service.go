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

package skeleton

import (
	"flag"
	"fmt"
	"github.com/houyi-tracing/houyi/pkg/hc"
	"github.com/houyi-tracing/houyi/pkg/skeleton/server"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	pMetrics "github.com/jaegertracing/jaeger/pkg/metrics"
)

type Service struct {
	ServiceName string

	Logger *zap.Logger

	AdminServer *server.AdminServer

	signalsChannel chan os.Signal

	hcStatusChannel chan hc.Status

	MetricsFactory metrics.Factory
}

func NewService(serviceName string, adminPort int) *Service {
	signalChannel := make(chan os.Signal, 1)
	hcStatusChannel := make(chan hc.Status)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	return &Service{
		ServiceName:     serviceName,
		AdminServer:     server.NewAdminServer(adminPort),
		Logger:          zap.NewNop(),
		signalsChannel:  signalChannel,
		hcStatusChannel: hcStatusChannel,
	}
}

func (s *Service) AddFlags(flagSet *flag.FlagSet) {
	s.AdminServer.AddFlags(flagSet)
	AddFlags(flagSet)
}

func (s *Service) SetHCStatus(status hc.Status) {
	s.hcStatusChannel <- status
}

func (s *Service) HC() hc.HealthCheck {
	return s.AdminServer.HC()
}

func (s *Service) Start(v *viper.Viper) error {
	sharedFlags := new(SharedFlags).InitFromViper(v)
	if logger, err := sharedFlags.NewLogger(zap.NewProductionConfig()); err == nil {
		s.Logger = logger
	} else {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	s.Logger.Info("Starting service", zap.String("service name", s.ServiceName))

	metricsBuilder := new(pMetrics.Builder).InitFromViper(v)
	metricsFactory, err := metricsBuilder.CreateMetricsFactory("")
	if err != nil {
		return fmt.Errorf("cannot create metrics factory: %w", err)
	}
	s.MetricsFactory = metricsFactory

	s.AdminServer.InitFromViper(v, s.Logger)
	if err := s.AdminServer.Serve(); err != nil {
		return fmt.Errorf("failed to start the admin server: %w", err)
	}
	return nil
}

func (s *Service) RunAndThen(shutdown func()) {
	s.HC().Ready()

statusLoop:
	for {
		select {
		case status := <-s.hcStatusChannel:
			s.HC().Set(status)
		case <-s.signalsChannel:
			break statusLoop
		}
	}

	s.Logger.Info("Shutting down...")
	s.HC().Set(hc.Unavailable)

	if shutdown != nil {
		shutdown()
	}

	if err := s.AdminServer.Close(); err != nil {
		s.Logger.Error("failed to close the admin server", zap.Error(err))
	}
	s.Logger.Info("Shutdown completed")
}
