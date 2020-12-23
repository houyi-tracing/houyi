// Copyright (c) 2020 The Houyi Authors.
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

package graph

import (
	"fmt"
	"time"
)

type expirationTimer map[ExecutionGraphNode]time.Time

func newExpirationTimer() nodeExpirationTimer {
	eT := make(expirationTimer)
	return &eT
}

func (e *expirationTimer) Timing(node ExecutionGraphNode, duration time.Duration) {
	(*e)[node] = time.Now().Add(duration)
}

func (e *expirationTimer) IsExpired(node ExecutionGraphNode) (bool, error) {
	if t, has := (*e)[node]; has {
		if t.Before(time.Now()) {
			delete(*e, node)
			return true, nil
		}
		return false, nil
	} else {
		return true, fmt.Errorf("node doest not exist in expiration timer: %v, %v", node.Service(), node.Operation())
	}
}

func (e *expirationTimer) Remove(node ExecutionGraphNode) {
	delete(*e, node)
}
