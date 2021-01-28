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
	uatomic "go.uber.org/atomic"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	minCapacity = 65536 // 2 ^ 16

	// When the queue need to resize, it would increase the size by multiplying 2 if previous capacity
	// is less than threshold, otherwise by multiplying by amplification.
	// This mechanism is similar to the append function of slices in Golang.
	threshold     = 1048576 // 2 ^ 20
	amplification = 1.25
	shrink        = 0.6

	resizeInterval = time.Second * 5
)

type DynamicQueue interface {
	Capacity() int
	Produce(item interface{}) bool
	Size() int
	StartConsumers(workers int, consumer func(item interface{}))
	Stop()
	resize()
}

type dq struct {
	workers int

	size *uatomic.Uint32

	capacity *uatomic.Uint32

	// use pointer(*) instead of reference for easier replacing queue in resizing
	items *chan interface{}

	consumer func(item interface{})

	resizeMu sync.RWMutex

	stopCh chan *sync.WaitGroup

	stopTickerCh chan *sync.WaitGroup
}

func NewDynamicQueue() DynamicQueue {
	queue := make(chan interface{}, minCapacity)
	q := &dq{
		size:         uatomic.NewUint32(0),
		capacity:     uatomic.NewUint32(minCapacity),
		items:        &queue,
		stopTickerCh: make(chan *sync.WaitGroup),
	}
	q.background(resizeInterval)
	return q
}

func (q *dq) Capacity() int {
	return int(q.capacity.Load())
}

func (q *dq) Produce(item interface{}) bool {
	q.size.Add(1)

	q.resizeMu.RLock()
	defer q.resizeMu.RUnlock()

	*q.items <- item
	return true
}

func (q *dq) Size() int {
	return int(q.size.Load())
}

func (q *dq) StartConsumers(workers int, consumer func(item interface{})) {
	q.workers = workers
	q.consumer = consumer
	q.stopCh = make(chan *sync.WaitGroup, workers)
	var startWg sync.WaitGroup

	for i := 0; i < q.workers; i++ {
		startWg.Add(1)
		go func() {
			startWg.Done()
			for {
				queue := *q.items
				select {
				case item, ok := <-queue:
					if ok {
						q.size.Sub(1)
						q.consumer(item)
					} else {
						// channel closed
						return
					}
				case wg := <-q.stopCh:
					// actively close whole queue
					wg.Done()
					return
				}
			}
		}()
	}

	startWg.Wait()
}

func (q *dq) Stop() {
	var wg sync.WaitGroup
	for i := 0; i < q.workers; i++ {
		wg.Add(1)
		q.stopCh <- &wg
	}
	wg.Add(1)
	q.stopTickerCh <- &wg
	wg.Wait() // wait for closing all consumers
	close(q.stopTickerCh)
	close(q.stopCh)
	close(*q.items)
}

func (q *dq) background(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				q.resize()
			case wg := <-q.stopTickerCh:
				wg.Done()
				return
			}
		}
	}()
}

func (q *dq) resize() {
	var newCapacity int

	size, capacity := q.Size(), q.Capacity()

	if size >= capacity {
		// increase capacity
		if capacity >= threshold {
			newCapacity = int(float64(capacity) * amplification)
		} else {
			newCapacity = capacity * 2
		}
	} else if size <= capacity/2 {
		// decrease capacity
		newCapacity = int(float64(capacity) * shrink)
		if newCapacity < minCapacity {
			return
		}
	} else {
		// do not need to resize when capacity/2 < size < capacity
		return
	}

	q.resizeMu.Lock()
	defer q.resizeMu.Unlock()

	previous := *q.items
	queue := make(chan interface{}, newCapacity)

	// replace old pointer of q.items with pointer of new queue.
	swapped := atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.items)),
		unsafe.Pointer(q.items),
		unsafe.Pointer(&queue))

	if swapped {
		// swapping the pointers of channels would not copy the data in these channels, therefore we should
		// send all data in previous channel to new one.
		q.StartConsumers(q.workers, q.consumer)

		close(previous)

		for i := 0; i < len(previous); i++ {
			if item, ok := <-previous; ok {
				queue <- item
			} else {
				break
			}
		}

		q.capacity.Store(uint32(newCapacity))
	}
}
