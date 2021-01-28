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

package app

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"math/rand"
	"sync"
	"time"
)

// SeedSet is a data structure to store routing information of peers(Seed) for Registry.
type SeedSet interface {
	Add(ip string, port int) int64
	AllPeers(self int64) []*api_v1.Peer
	GetNode(id int64) *api_v1.Peer
	Has(id int64) bool
	IsDead(id int64, life time.Duration) bool
	Remove(id int64)
	Refresh(id int64)
	Update(id int64, ip string, port int)
	AllIds() []int64
}

type item struct {
	upSince time.Time
	seed    *api_v1.Peer
}

type seedSet struct {
	lock *sync.RWMutex
	m    map[int64]*item
}

func NewSeedSet() SeedSet {
	return &seedSet{
		lock: &sync.RWMutex{},
		m:    make(map[int64]*item),
	}
}

func (s *seedSet) Add(ip string, port int) int64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	newId := s.newId()
	s.m[newId] = &item{
		upSince: time.Now(),
		seed: &api_v1.Peer{
			Ip:   ip,
			Port: int64(port),
		},
	}
	return newId
}

func (s *seedSet) AllPeers(self int64) []*api_v1.Peer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	ret := make([]*api_v1.Peer, 0)
	for id, t := range s.m {
		if id != self {
			ret = append(ret, t.seed)
		}
	}
	return ret
}

func (s *seedSet) AllIds() []int64 {
	ret := make([]int64, 0, len(s.m))
	for k := range s.m {
		ret = append(ret, k)
	}
	return ret
}

func (s *seedSet) Has(id int64) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, has := s.m[id]
	return has
}

func (s *seedSet) GetNode(id int64) *api_v1.Peer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if i, has := s.m[id]; has {
		return i.seed
	} else {
		return nil
	}
}

func (s *seedSet) IsDead(id int64, life time.Duration) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if i, has := s.m[id]; has && i.upSince.Add(life).Before(time.Now()) {
		return true
	} else {
		return false
	}
}

func (s *seedSet) Remove(id int64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.m, id)
}

func (s *seedSet) Refresh(id int64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if i, has := s.m[id]; has {
		i.upSince = time.Now()
	}
}

func (s *seedSet) Update(id int64, ip string, port int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if i, has := s.m[id]; has {
		i.seed.Ip = ip
		i.seed.Port = int64(port)
	}
}

func (s *seedSet) newId() int64 {
	var id int64
	rand.Seed(time.Now().Unix())
	for {
		id = rand.Int63n(maxNode)
		if _, has := s.m[id]; !has {
			return id
		}
	}
}
