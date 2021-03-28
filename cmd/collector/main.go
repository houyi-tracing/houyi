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

package main

import (
	"fmt"
	"github.com/houyi-tracing/houyi/cmd/collector/app"
	"github.com/houyi-tracing/houyi/cmd/collector/app/filter"
	"github.com/houyi-tracing/houyi/cmd/collector/app/processor"
	"github.com/houyi-tracing/houyi/pkg/config"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip/seed"
	"github.com/houyi-tracing/houyi/pkg/gossip/server"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/skeleton"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/jaegertracing/jaeger/plugin/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"log"
	"os"
)

const (
	serviceName = "houyi-collector"
)

func main() {
	v := viper.New()
	v.AutomaticEnv() // read env params.

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	svc := skeleton.NewService(serviceName, ports.AdminHttpPort)

	storageFactory, err := storage.NewFactory(storage.FactoryConfigFromEnvAndCLI(os.Args, os.Stderr))
	if err != nil {
		log.Fatalf("Cannot initialize storage factory: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   serviceName,
		Short: "Collector for Houyi tracing",
		Long:  `This is a collector to receive and process spans reported by agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}

			logger := svc.Logger // for short

			// Trace Graph
			traceGraph := tg.NewTraceGraph(logger)

			// evaluator
			eval := evaluator.NewEvaluator(logger)

			// Filter
			sf := filter.NewSpanFilter()

			// Gossip Seed
			logger.Info("Starting gossip seed")
			seedOpts := new(seed.Flags).InitFromViper(v)
			gossipSeed, err := server.BuildSeed(&server.SeedParams{
				Logger:     logger,
				ListenPort: seedOpts.SeedGrpcPort,
				LruSize:    seedOpts.LruSize,
				ConfigServerEndpoint: &routing.Endpoint{
					Addr: seedOpts.ConfigServerAddress,
					Port: seedOpts.ConfigServerGrpcPort,
				},
				TraceGraph: traceGraph,
				Evaluator:  eval,
			})
			if err != nil {
				return err
			}

			// reuse span writer of Jaeger
			baseFactory := svc.MetricsFactory.Namespace(metrics.NSOptions{Name: "houyi"})
			storageFactory.InitFromViper(v)
			if err := storageFactory.Initialize(baseFactory, logger); err != nil {
				logger.Fatal("Failed to init storage factory", zap.Error(err))
				return err
			}
			sw, err := storageFactory.CreateSpanWriter()
			if err != nil {
				logger.Fatal("Failed to create span writer", zap.Error(err))
				return err
			}

			// Span Processor
			logger.Info("Initializing span processor")
			spOpts := new(processor.Flags).InitFromViper(v)
			sp := processor.NewSpanProcessor(logger,
				processor.Options.NumWorkers(spOpts.NumWorkers),
				processor.Options.GossipSeed(gossipSeed),
				processor.Options.TraceGraph(traceGraph),
				processor.Options.EvaluateSpan(eval.Evaluate),
				processor.Options.FilterSpan(sf.Filter),
				processor.Options.SpanWriter(sw),
				processor.Options.ConfigServerEndpoint(&routing.Endpoint{
					Addr: spOpts.ConfigServerAddr,
					Port: spOpts.ConfigServerPort,
				}))

			// Collector
			cOpts := new(app.Flags).InitFromViper(v)
			c := app.NewCollector(&app.CollectorParams{
				Logger:         logger,
				SpanProcessor:  sp,
				GrpcListenPort: cOpts.GrpcListenPort,
				Evaluator:      eval,
			})

			if err = gossipSeed.Start(); err != nil {
				logger.Fatal("failed to start gossip seed", zap.Error(err))
				return err
			} else {
				logger.Info("Started gossip seed")
			}

			if err = c.Start(); err != nil {
				logger.Fatal("Failed to start collector", zap.Error(err))
				return err
			} else {
				logger.Info("Started collector")
			}

			svc.RunAndThen(func() {
				// Do some nothing before completing shutting down.
				// for example, closing I/O or DB connection, etc.
				if err := c.Close(); err != nil {
					logger.Fatal("Failed to close collector", zap.Error(err))
				}
				if err := gossipSeed.Stop(); err != nil {
					logger.Fatal("Failed to stop gossip seed", zap.Error(err))
				}
			})
			return nil
		},
	}

	config.AddFlags(
		v,
		rootCmd,
		processor.AddFlags,
		seed.AddFlags,
		app.AddFlags,
		storageFactory.AddFlags,
		svc.AddFlags)

	// rootCmd represents the base command when called without any subcommands
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
