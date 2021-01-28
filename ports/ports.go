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

package ports

const (
	// Admin Server
	AdminHttpPort = 22590 // port to serve HTTP for health checking

	// Gossip
	RegistryGrpcListenPort = 22600 // port to serve gRPC for Gossip Registry
	SeedGrpcListenPort     = 22650 // port to serve gRPC for Gossip Seed

	// Strategy Manager
	StrategyManagerGrpcListenPort = 18760 // port to server gRPC for StrategyManager.

	// Collector
	CollectorGrpcListenPort = 14580

	// Agent
	AgentGrpcListenPort = 14680
)
