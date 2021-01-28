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
	"time"
)

type options struct {
	randomPick      int
	listenPort      int
	probToR         float64
	refreshInterval time.Duration
}

type Option func(o *options)

var Options options

func (options) RandomPick(n int) Option {
	return func(opts *options) {
		opts.randomPick = n
	}
}

func (options) ListenPort(p int) Option {
	return func(opts *options) {
		opts.listenPort = p
	}
}

func (options) ProbToR(p float64) Option {
	return func(opts *options) {
		opts.probToR = p
	}
}

func (options) RefreshInterval(d time.Duration) Option {
	return func(o *options) {
		o.refreshInterval = d
	}
}

func (o *options) apply(opts ...Option) options {
	for _, opt := range opts {
		opt(o)
	}
	if o.randomPick == 0 {
		o.randomPick = DefaultRandomPick
	}
	if o.listenPort == 0 {
		o.listenPort = DefaultListenPort
	}
	if o.probToR == 0 {
		o.probToR = DefaultProbToR
	}
	if o.refreshInterval == 0 {
		o.refreshInterval = DefaultRefreshInterval
	}
	return *o
}
