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

// Sampling strategy tree.
package sst

import "github.com/houyi-tracing/houyi/idl/api_v1"

type SamplingStrategyTree interface {
	// Add adds a new operation into this tree.
	Add(op *api_v1.Operation) error

	// Has returns true if operation already exists in this tree, else false.
	Has(op *api_v1.Operation) bool

	// Generate generates sampling strategy for inputted operation.
	Generate(op *api_v1.Operation) (float64, error)

	// Promote promotes a operation in this tree. As a result, the sampling rate of inputted operation will increase.
	Promote(op *api_v1.Operation) error

	// Prune removes inputted operation from this tree.
	Prune(op *api_v1.Operation) error
}
