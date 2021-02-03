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
	"github.com/houyi-tracing/houyi/cmd/sm/app"
	"github.com/houyi-tracing/houyi/pkg/config"
	"github.com/houyi-tracing/houyi/pkg/evaluator"
	"github.com/houyi-tracing/houyi/pkg/gossip/seed"
	"github.com/houyi-tracing/houyi/pkg/gossip/server"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/skeleton"
	"github.com/houyi-tracing/houyi/pkg/sst"
	"github.com/houyi-tracing/houyi/pkg/tg"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"os"
)

const (
	serviceName = "strategy-manger"
)

func main() {
	v := viper.New()
	v.AutomaticEnv() // read env params.

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}
	svc := skeleton.NewService(serviceName, ports.AdminHttpPort)

	rootCmd := &cobra.Command{
		Use:   serviceName,
		Short: "Strategy manager for Houyi tracing",
		Long: `This is a server to store sampling strategies and process requests for 
			pulling sampling strategies from agents.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}

			logger := svc.Logger // for short

			// Sampling Strategy Tree
			sstOpts := new(sst.Flags).InitFromViper(v)
			strategyStore := sst.NewSamplingStrategyTree(sstOpts.MaxChildNodes)

			// Trace Graph
			traceGraph := tg.NewTraceGraph(logger)

			// Evaluator
			eval := evaluator.NewEvaluator(logger)

			// Gossip Seed
			seedOpts := new(seed.Flags).InitFromViper(v)
			gossipSeed, err := server.CreateAndStartSeed(&server.SeedParams{
				Logger:     logger,
				ListenPort: seedOpts.SeedGrpcPort,
				LruSize:    seedOpts.LruSize,
				RegistryEndpoint: &routing.Endpoint{
					Addr: seedOpts.RegistryAddress,
					Port: seedOpts.RegistryGrpcPort,
				},
				TraceGraph: traceGraph,
				Evaluator:  eval,
			})
			if err != nil {
				return err
			}

			// Strategy Manager
			smOpts := new(app.Flags).InitFromViper(v)
			sm := app.NewStrategyManager(&app.StrategyManagerParams{
				Logger:          logger,
				GrpcListenPort:  smOpts.GrpcListenPort,
				RefreshInterval: smOpts.RefreshInterval,
				StrategyStore:   strategyStore,
				TraceGraph:      traceGraph,
				Evaluator:       eval,
				GossipSeed:      gossipSeed,
				ScaleFactor:     atomic.NewFloat64(smOpts.ScaleFactor),
			})
			if err := sm.Start(); err != nil {
				return err
			}

			svc.RunAndThen(func() {
				// Do something before completing shutting down.
				// for example, closing I/O or DB connection, etc.
				if err := sm.Stop(); err != nil {
					logger.Error("Failed to stop strategy manager", zap.Error(err))
				}
			})
			return nil
		},
	}

	config.AddFlags(
		v,
		rootCmd,
		seed.AddFlags,
		sst.AddFlags,
		app.AddFlags,
		svc.AddFlags)

	// rootCmd represents the base command when called without any subcommands
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
