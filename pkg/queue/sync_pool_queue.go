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
	"fmt"
	"math"
	"sync"
	"time"
)

type qItem struct {
	next *qItem
	val  interface{}
}

type syncPoolQueue struct {
	sync.RWMutex

	pool     sync.Pool
	size     int
	capacity int
	workers  int

	head *qItem
	tail *qItem

	stopCh chan *sync.WaitGroup
}

func NewSyncPoolQueue(capacity int) DynamicQueue {
	head, tail := &qItem{}, &qItem{}
	head.next = tail
	q := &syncPoolQueue{
		size:     0,
		capacity: capacity,
		head:     head,
		tail:     tail,
		pool: sync.Pool{
			New: func() interface{} {
				return &qItem{}
			},
		},
		stopCh: make(chan *sync.WaitGroup),
	}
	return q
}

func (q *syncPoolQueue) Capacity() int {
	return math.MaxInt64
}

func (q *syncPoolQueue) Produce(item interface{}) bool {
	return q.pushBack(item)
}

func (q *syncPoolQueue) Size() int {
	q.RLock()
	defer q.RUnlock()
	return q.size
}

func (q *syncPoolQueue) StartConsumers(workers int, consumer func(item interface{})) {
	q.workers = workers
	for i := 0; i < q.workers; i++ {
		go func() {
			for {
				select {
				case wg := <-q.stopCh:
					wg.Done()
					return
				default:
					if item, err := q.popFront(); err == nil {
						consumer(item)
					} else {
						time.Sleep(time.Second) // sleep one second if queue is empty
					}
				}
			}
		}()
	}
}

func (q *syncPoolQueue) Stop() {
	var wg sync.WaitGroup
	wg.Add(q.workers)
	for i := 0; i < q.workers; i++ {
		q.stopCh <- &wg
	}
	wg.Wait()
	close(q.stopCh)
}

func (q *syncPoolQueue) pushBack(val interface{}) bool {
	q.Lock()
	defer q.Unlock()

	if q.size < q.capacity {
		qi := q.pool.Get().(*qItem)
		qi.val = val

		q.tail.next = qi
		q.tail = qi
		q.size += 1
		return true
	} else {
		return false
	}
}

func (q *syncPoolQueue) popFront() (interface{}, error) {
	q.Lock()
	defer q.Unlock()

	if q.head.next != q.tail {
		first := q.head.next
		retVal := first.val

		// head -> first -> second => head -> second
		q.head.next = first.next

		q.pool.Put(first) // reuse memory
		q.size -= 1
		return retVal, nil
	} else {
		return nil, fmt.Errorf("can not do pop operation on an empty queue")
	}
}
