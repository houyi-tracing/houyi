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

package tg

import "github.com/houyi-tracing/houyi/idl/api_v1"

type TraceNode struct {
	Name     string       `json:"name"`
	Children []*TraceNode `json:"children"`
}

type TraceGraph interface {
	// Add adds an operation into trace graph and returned error would be nil if the operation did not exist in advance.
	Add(op *api_v1.Operation) error

	// Remove removes an operation from trace graph and returned error would be nil if the operation already exist.
	Remove(op *api_v1.Operation) error

	// Has returns true if operation exist in trace graph, else false.
	Has(op *api_v1.Operation) bool

	// AddRelation adds a new relation in trace graph. The error it returns would be nil if both operations already
	// exist in trace graph.
	AddRelation(rel *api_v1.Relation) error

	// RemoveRelation removes old relation in trace graph. The error it returns would be nil if both operations already
	// exist in trace graph.
	RemoveRelation(rel *api_v1.Relation) error

	// HasRelation returns true if there is a relation between inputted operations, else false.
	// If one of operations of this relation does not exist, this function would return false.
	HasRelation(rel *api_v1.Relation) bool

	// IsIngress returns true if operation is an ingress operation that received requests from users instead of other
	// operations in this application system.
	IsIngress(op *api_v1.Operation) bool

	// GetIngresses returns all ingress operations relate to inputted operation.
	GetIngresses(op *api_v1.Operation) ([]*api_v1.Operation, error)

	// AllIngresses returns all ingress operations.
	AllIngresses() []*api_v1.Operation

	// Traces returns static call relationships relate to inputted operation.
	// Each element of returned slice is an entry operation.
	Dependencies(op *api_v1.Operation) ([]*TraceNode, error)

	Services() []string

	Operations(string) []string

	// Size returns the number of operations in this trace graph.
	Size() int
}
