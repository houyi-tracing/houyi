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

package queue

import (
	"github.com/bmizerany/assert"
	"go.uber.org/atomic"
	"testing"
	"time"
)

func TestSyncPoolQueue(t *testing.T) {
	maxCapacity := 10000000
	n := 100000

	q := NewSyncPoolQueue(maxCapacity)
	i := atomic.Int64{}
	q.StartConsumers(10, func(item interface{}) {
		i.Add(1)
	})

	for i := 0; i < n; i++ {
		q.Produce(i)
	}

	time.Sleep(time.Second * 5)

	q.Stop()

	assert.Equal(t, int64(n), i.Load())
}
