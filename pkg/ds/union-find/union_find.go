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
	"fmt"
	"log"
	"sync"
)

// unionFind is a data structure that stores a collection of disjoint (non-overlapping) sets.
type unionFind struct {
	m  map[interface{}]interface{}
	mu *sync.RWMutex
}

func NewUnionFind() UnionFind {
	return &unionFind{
		m:  make(map[interface{}]interface{}),
		mu: new(sync.RWMutex),
	}
}

// GetAll returns a slice contains all elements in union-find.
func (uf *unionFind) GetAll() []interface{} {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	all := make([]interface{}, 0, len(uf.m))
	for k := range uf.m {
		all = append(all, k)
	}
	return all
}

// Add adds a new element into this union-find as a disjoint set.
// If the element is already in union-find, this operation would be ignored.
func (uf *unionFind) Add(v interface{}) {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	if _, has := uf.m[v]; !has {
		uf.m[v] = v
	}
}

// Find returns the representative of the set v belongs to.
func (uf *unionFind) Find(v interface{}) interface{} {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	if ret, err := uf.find(v); err != nil {
		log.Println(err)
		return nil
	} else {
		return ret
	}
}

// Has return true if v exist in union-find, else false.
func (uf *unionFind) Has(v interface{}) bool {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	if _, has := uf.m[v]; has {
		return true
	} else {
		return false
	}
}

// Union merges the two sets containing a and b, respectively.
func (uf *unionFind) Union(v1 interface{}, v2 interface{}) error {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	if _, has := uf.m[v1]; !has {
		return fmt.Errorf("can not union non-existent elements: %v", v1)
	}
	if _, has := uf.m[v2]; !has {
		return fmt.Errorf("can not union non-existent elements: %v", v2)
	}

	p1, _ := uf.find(v1)
	p2, _ := uf.find(v2)
	uf.m[p1] = p2
	return nil
}

// Size returns the number of elements in union-find.
func (uf *unionFind) Size() int {
	uf.mu.RLock()
	defer uf.mu.RUnlock()

	return len(uf.m)
}

// Remove removes an exist element in union-find.
// If the element to be removed is the representative of a set, all elements in that set would be removed.
func (uf *unionFind) Remove(toRemoved interface{}) {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	if _, has := uf.m[toRemoved]; has {
		for k := range uf.m {
			if find, err := uf.find(k); err != nil || find == toRemoved {
				delete(uf.m, k)
			}
		}
	}
}

func (uf *unionFind) find(v interface{}) (interface{}, error) {
	if parent, has := uf.m[v]; has {
		if parent != v {
			uf.m[v], _ = uf.find(parent) // path compression
		}
		return uf.m[v], nil
	} else {
		return nil, fmt.Errorf("union-find set has not element: %v", v)
	}
}
