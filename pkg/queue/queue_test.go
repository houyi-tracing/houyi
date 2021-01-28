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
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

type queueItem struct {
	val int
}

func TestSwapChanPointer(t *testing.T) {
	a := make(chan int, 5)
	ap := &a

	for i := 0; i < 2; i++ {
		*ap <- i + 100
	}

	b := make(chan int, 10)
	bp := &b

	fmt.Printf("a\t%p\n", a)
	fmt.Printf("b\t%p\n", b)
	fmt.Printf("*ap\t%p\n", *ap)
	fmt.Printf("ap\t%v\n", ap)
	fmt.Printf("*bp\t%p\n", *bp)
	fmt.Printf("bp\t%v\n", bp)

	swapped := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ap)), unsafe.Pointer(ap), unsafe.Pointer(bp))

	if swapped {
		fmt.Println()
		fmt.Printf("a\t%p\n", a)
		fmt.Printf("b\t%p\n", b)
		fmt.Printf("*ap\t%p\n", *ap)
		fmt.Printf("ap\t%v\n", ap)
		fmt.Printf("*bp\t%p\n", *bp)
		fmt.Printf("bp\t%v\n", bp)
	}
}

func TestResizeWhenSizeIsGreaterThanOrEqualToCapacity(t *testing.T) {
	cpuN := 4
	runtime.GOMAXPROCS(cpuN)

	q := NewDynamicQueue()
	l := &sync.Mutex{}
	cnt := 0
	q.StartConsumers(cpuN, func(item interface{}) {
		l.Lock()
		defer l.Unlock()
		time.Sleep(time.Microsecond)
		cnt += 1
	})
	fmt.Println("Started consumers")

	startTime := time.Now()

	n := 10000000
	for i := 0; i < n; i++ {
		q.Produce(queueItem{val: i})
	}
	fmt.Println("Finished producing items")

	for {
		time.Sleep(time.Second)
		if cnt >= n {
			fmt.Println("Consumed all items")
			break
		}
	}

	fmt.Println(time.Now().Sub(startTime).String())
	assert.Equal(t, 0, q.Size())
	assert.Equal(t, n, cnt)
}
