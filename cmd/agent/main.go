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
	"github.com/houyi-tracing/houyi/cmd/agent/app"
	"github.com/houyi-tracing/houyi/pkg/config"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/pkg/skeleton"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const (
	serviceName = "houyi-agent"
)

func main() {
	v := viper.New()
	v.AutomaticEnv() // read env params.

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	svc := skeleton.NewService(serviceName, ports.AdminHttpPort)

	var rootCmd = &cobra.Command{
		Use:   serviceName,
		Short: "Collector for Houyi tracing",
		Long:  `This is a collector to receive and process spans reported by agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}

			logger := svc.Logger

			aOpts := new(app.Flags).InitFromViper(v)
			a := app.NewAgent(&app.AgentParams{
				Logger:         logger,
				GrpcListenPort: aOpts.GrpcListenPort,
				CollectorEndpoint: &routing.Endpoint{
					Addr: aOpts.CollectorAddr,
					Port: aOpts.CollectorPort,
				},
				ConfigServerEp: &routing.Endpoint{
					Addr: aOpts.ConfigServerAddr,
					Port: aOpts.ConfigServerPort,
				},
			})

			if err := a.Start(); err != nil {
				return err
			}

			svc.RunAndThen(func() {
				// Do some nothing before completing shutting down.
				// for example, closing I/O or DB connection, etc.
				if err := a.Stop(); err != nil {
					logger.Fatal("Failed to stop agent", zap.Error(err))
				}
			})
			return nil
		},
	}

	config.AddFlags(
		v,
		rootCmd,
		app.AddFlags,
		svc.AddFlags)

	// rootCmd represents the base command when called without any subcommands
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
