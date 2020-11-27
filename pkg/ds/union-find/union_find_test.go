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

package union_find

import (
	"github.com/bmizerany/assert"
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	uf := NewUnionFind()
	for _, n := range []int{1, 2, 3, 4, 5} {
		uf.Add(n)
	}
	assert.Equal(t, 5, len(uf.GetAll()))
}

func TestFind(t *testing.T) {
	uf := NewUnionFind()
	for _, n := range []int{1, 2, 3, 4, 5} {
		uf.Add(n)
	}

	for i := 1; i <= 5; i++ {
		find := uf.Find(i)
		assert.Equal(t, i, find)
	}
}

func TestUnion(t *testing.T) {
	uf := NewUnionFind()
	for _, n := range []int{1, 2, 3, 4, 5} {
		uf.Add(n)
	}

	for i := 1; i < 5; i++ {
		_ = uf.Union(i, 5)
	}

	all := uf.GetAll()
	for _, n := range all {
		find := uf.Find(n)
		assert.Equal(t, find, 5)
	}
}

func TestRemove(t *testing.T) {
	uf := NewUnionFind()
	for _, n := range []int{1, 2, 3, 4, 5} {
		uf.Add(n)
	}

	_ = uf.Union(1, 3)
	assert.Equal(t, 3, uf.Find(1))

	_ = uf.Union(3, 5)
	assert.Equal(t, 5, uf.Find(3))

	_ = uf.Union(5, 4)
	assert.Equal(t, 4, uf.Find(5))

	uf.Remove(4)
	all := uf.GetAll()
	assert.Equal(t, 1, len(all))
}

func TestHas(t *testing.T) {
	uf := NewUnionFind()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(uf UnionFind, wg *sync.WaitGroup) {
			defer wg.Done()

			uf.Has(1)
		}(uf, &wg)
	}

	wg.Wait()
}
