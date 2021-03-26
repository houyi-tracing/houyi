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

package store

import (
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"sync"
)

type StrategyStore interface {
	Add(svc, op string, strategy *api_v1.PerOperationStrategy)
	Update(svc, op string, strategy *api_v1.PerOperationStrategy)
	UpdateAll([]*api_v1.PerOperationStrategy)
	Has(svc, op string) bool
	Override([]*api_v1.PerOperationStrategy)
	Get(svc, op string) (*api_v1.PerOperationStrategy, error)
	GetAll() []*api_v1.PerOperationStrategy
	SetDefaultStrategy(*api_v1.PerOperationStrategy)
	GetDefaultStrategy() *api_v1.PerOperationStrategy
	Remove(svc, op string) error
	RemoveAll()
}

const (
	operationNotExist = "operation does not exist"
)

type strategyStore struct {
	sync.RWMutex

	defaultStrategy *api_v1.PerOperationStrategy
	data            map[string]map[string]*api_v1.PerOperationStrategy
}

func NewStrategyStore() StrategyStore {
	return &strategyStore{
		defaultStrategy: &api_v1.PerOperationStrategy{
			Type: api_v1.Type_DYNAMIC,
			Strategy: &api_v1.PerOperationStrategy_Dynamic{
				Dynamic: &api_v1.DynamicSampling{
					SamplingRate: 1.0,
				},
			},
		},
		data: make(map[string]map[string]*api_v1.PerOperationStrategy),
	}
}

func (store *strategyStore) Add(svc, op string, strategy *api_v1.PerOperationStrategy) {
	store.Lock()
	defer store.Unlock()

	store.add(svc, op, strategy)
}

func (store *strategyStore) Update(svc, op string, strategy *api_v1.PerOperationStrategy) {
	store.Lock()
	defer store.Unlock()

	store.add(svc, op, strategy)
}

func (store *strategyStore) UpdateAll(strategies []*api_v1.PerOperationStrategy) {
	store.Lock()
	defer store.Unlock()

	for _, s := range strategies {
		store.add(s.GetService(), s.GetOperation(), s)
	}
}

func (store *strategyStore) SetDefaultStrategy(strategy *api_v1.PerOperationStrategy) {
	store.Lock()
	defer store.Unlock()

	store.defaultStrategy = strategy
}

func (store *strategyStore) GetDefaultStrategy() *api_v1.PerOperationStrategy {
	store.RLock()
	defer store.RUnlock()

	return store.defaultStrategy
}

func (store *strategyStore) Has(svc, op string) bool {
	store.RLock()
	defer store.RUnlock()

	return store.has(svc, op)
}

func (store *strategyStore) Override(strategies []*api_v1.PerOperationStrategy) {
	store.Lock()
	defer store.Unlock()

	store.data = make(map[string]map[string]*api_v1.PerOperationStrategy)
	for _, s := range strategies {
		store.add(s.Service, s.Operation, s)
	}
}

func (store *strategyStore) Get(svc, op string) (*api_v1.PerOperationStrategy, error) {
	store.RLock()
	defer store.RUnlock()

	if store.has(svc, op) {
		return store.data[svc][op], nil
	} else {
		return nil, fmt.Errorf(operationNotExist)
	}
}

func (store *strategyStore) GetAll() []*api_v1.PerOperationStrategy {
	store.RLock()
	defer store.RUnlock()

	ret := make([]*api_v1.PerOperationStrategy, 0)
	for _, opMap := range store.data {
		for _, pos := range opMap {
			ret = append(ret, pos)
		}
	}
	return ret
}

func (store *strategyStore) Remove(svc, op string) error {
	store.Lock()
	defer store.Unlock()

	if store.has(svc, op) {
		delete(store.data[svc], op)
		return nil
	} else {
		return fmt.Errorf(operationNotExist)
	}
}

func (store *strategyStore) RemoveAll() {
	store.Lock()
	defer store.Unlock()

	store.data = make(map[string]map[string]*api_v1.PerOperationStrategy)
}

func (store *strategyStore) has(svc, op string) bool {
	if _, ok := store.data[svc]; ok {
		if _, has := store.data[svc][op]; has {
			return true
		}
	}
	return false
}

func (store *strategyStore) add(svc, op string, strategy *api_v1.PerOperationStrategy) {
	if _, ok := store.data[svc]; !ok {
		store.data[svc] = make(map[string]*api_v1.PerOperationStrategy)
	}
	store.data[svc][op] = strategy
}
