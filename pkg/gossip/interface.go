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

package gossip

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"time"
)

type Runnable interface {
	Start() error
	Stop() error
}

// Registry is used to help seeds(Seed) discover other peer nodes, so seeds no longer need to set up routing information
// for other peers in advance, which also provides runtime scalability.
type Registry interface {
	Runnable

	// Register receives register request from seeds at the initial phase of seeds and stores routing information
	// provided by the request.
	Register(ip string, port int) (int64, int, time.Duration, float64)

	// Heartbeat receives heartbeats from seeds and removes seeds that have not sent heartbeat messages for a long time.
	Heartbeat(nodeId int64, ip string, port int) (int64, []*api_v1.Peer)

	// AllPeers returns all alive peers
	AllSeeds() []*api_v1.Peer
}

type Seed interface {
	Runnable

	// OnNewRelation sets function that would be invoked when gossip seed received a message carrying new relation
	// and process it.
	OnNewRelation(func(rel *api_v1.Relation))

	// OnNewOperation sets function that would be invoked when gossip seed received a message carrying new operation
	// and process it.
	OnNewOperation(func(op *api_v1.Operation))

	// OnExpiredOperation sets function that would be invoked when gossip seed received a message carrying expired
	// operation and process it.
	OnExpiredOperation(func(op *api_v1.Operation))

	// OnEvaluatingTags sets function that would be invoked when gossip seed received a message carrying evaluating tags
	// process it.
	OnEvaluatingTags(func(tags *api_v1.EvaluatingTags))

	// MongerNewRelation activates message mongering to synchronize new relations between gossip seeds.
	MongerNewRelation(rel *api_v1.Relation)

	// MongerNewOperation activates message mongering to synchronize new operations between gossip seeds.
	MongerNewOperation(op *api_v1.Operation)

	// MongerExpiredOperation activates message mongering to synchronize expired operations between gossip seeds.
	MongerExpiredOperation(op *api_v1.Operation)

	// MongerEvaluatingTags activates message mongering to synchronize evaluating tags between gossip seeds.
	MongerEvaluatingTags(tags *api_v1.EvaluatingTags)
}
