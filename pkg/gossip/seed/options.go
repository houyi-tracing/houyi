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

package seed

import (
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/routing"
)

type options struct {
	lruSize            int
	listenPort         int
	registryEndpoint   *routing.Endpoint
	onNewRelation      func(rel *api_v1.Relation)
	onExpiredOperation func(op *api_v1.Operation)
	onEvaluatingTags   func(tags *api_v1.EvaluatingTags)
}

type Option func(opts *options)

var Options options

func (options) LruSize(s int) Option {
	return func(opts *options) {
		opts.lruSize = s
	}
}

func (options) ListenPort(p int) Option {
	return func(opts *options) {
		opts.listenPort = p
	}
}

func (options) RegistryEndpoint(ep *routing.Endpoint) Option {
	return func(opts *options) {
		opts.registryEndpoint = ep
	}
}

func (options) OnNewRelation(f func(rel *api_v1.Relation)) Option {
	return func(opts *options) {
		opts.onNewRelation = f
	}
}

func (options) OnExpiredOperation(f func(op *api_v1.Operation)) Option {
	return func(opts *options) {
		opts.onExpiredOperation = f
	}
}

func (options) OnFilterTags(f func(tags *api_v1.EvaluatingTags)) Option {
	return func(opts *options) {
		opts.onEvaluatingTags = f
	}
}

func (o options) apply(opts ...Option) options {
	ret := options{}
	for _, op := range opts {
		op(&ret)
	}

	if ret.registryEndpoint == nil {
		panic("GossipRegistry endpoint for gossip seed must not be empty")
	}
	if ret.lruSize <= 0 {
		ret.lruSize = DefaultLruSize
	}
	if ret.listenPort == 0 {
		ret.listenPort = DefaultSeedGrpcPort
	}
	if ret.onNewRelation == nil {
		ret.onNewRelation = func(rel *api_v1.Relation) {
			// do nothing
		}
	}
	if ret.onExpiredOperation == nil {
		ret.onExpiredOperation = func(op *api_v1.Operation) {
			// do nothing
		}
	}
	if ret.onEvaluatingTags == nil {
		ret.onEvaluatingTags = func(tags *api_v1.EvaluatingTags) {
			// do nothing
		}
	}

	return ret
}
