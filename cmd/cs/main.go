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
	"github.com/houyi-tracing/houyi/cmd/cs/app"
	"github.com/houyi-tracing/houyi/cmd/cs/app/registry"
	"github.com/houyi-tracing/houyi/cmd/cs/app/store"
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
	"go.uber.org/zap"
	"os"
)

const (
	serviceName = "central-configuration-server"
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
		Short: "Central configuration server of Houyi tracing",
		Long: `This is a server to store sampling strategies and evaluator.proto, and process requests for 
			pulling sampling strategies from agents and requests for promote operations from collectors.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}

			logger := svc.Logger // for short

			// Sampling Strategy Tree
			sstOpts := new(sst.Flags).InitFromViper(v)
			ssTree := sst.NewSamplingStrategyTree(sstOpts.Order)

			// Trace Graph
			traceGraph := tg.NewTraceGraph(logger)

			// evaluator
			eval := evaluator.NewEvaluator(logger)

			csOpts := new(app.Flags).InitFromViper(v)
			gossipRegistry := registry.NewRegistry(logger,
				csOpts.RandomPick,
				csOpts.ProbToR,
				csOpts.HeartbeatInterval)
			if err := gossipRegistry.Start(); err != nil {
				logger.Fatal("failed to start registry", zap.Error(err))
			}

			strategyStore := store.NewStrategyStore()

			// Gossip Seed
			seedOpts := new(seed.Flags).InitFromViper(v)
			gossipSeed, err := server.BuildSeed(&server.SeedParams{
				Logger:     logger,
				ListenPort: seedOpts.SeedGrpcPort,
				LruSize:    seedOpts.LruSize,
				ConfigServerEndpoint: &routing.Endpoint{
					Addr: "localhost",
					Port: seedOpts.ConfigServerGrpcPort,
				},
				TraceGraph: traceGraph,
				Evaluator:  eval,
			})
			if err != nil {
				return err
			}

			operationStore := store.NewOperationStore(logger, csOpts.OperationExpire, gossipSeed, ssTree, traceGraph)

			cs := app.NewConfigServer(&app.ConfigurationServerParams{
				Logger:          logger,
				GrpcListenPort:  csOpts.GrpcListenPort,
				HttpListenPort:  csOpts.HttpListenPort,
				GossipSeed:      gossipSeed,
				GossipRegistry:  gossipRegistry,
				TraceGraph:      traceGraph,
				Evaluator:       eval,
				StrategyStore:   strategyStore,
				SST:             ssTree,
				OperationStore:  operationStore,
				ScaleFactor:     csOpts.ScaleFactor,
				MinSamplingRate: csOpts.MinSamplingRate,
			})

			if err = cs.Start(); err != nil {
				return err
			}

			svc.RunAndThen(func() {
				// Do something before completing shutting down.
				// for example, closing I/O or DB connection, etc.
				if err = gossipRegistry.Stop(); err != nil {
					logger.Error("failed to stop registry", zap.Error(err))
				}
				if err = cs.Stop(); err != nil {
					logger.Error("failed to stop configuration server", zap.Error(err))
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
