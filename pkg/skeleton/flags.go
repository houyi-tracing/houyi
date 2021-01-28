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

package skeleton

import (
	"flag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevel = "log.level"
)

const (
	DefaultLogLevel = "info"
)

type SharedFlags struct {
	Logging
}

type Logging struct {
	LogLevel string
}

func AddFlags(flags *flag.FlagSet) {
	addLoggingFlags(flags)
}

func (aOpts *SharedFlags) InitFromViper(v *viper.Viper) *SharedFlags {
	aOpts.LogLevel = v.GetString(logLevel)
	return aOpts
}

func (aOpts *SharedFlags) NewLogger(conf zap.Config, options ...zap.Option) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(aOpts.LogLevel))
	if err != nil {
		return nil, err
	}
	conf.Level = zap.NewAtomicLevelAt(level)
	conf.Encoding = "console"
	conf.EncoderConfig = zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "Time",
		LevelKey:       "Level",
		NameKey:        "Name",
		CallerKey:      "Caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "Message",
		StacktraceKey:  "Stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return conf.Build(options...)
}

func addLoggingFlags(flags *flag.FlagSet) {
	flags.String(logLevel, DefaultLogLevel,
		"Minimal allowed log Level. For more levels see https://github.com/uber-go/zap")
}
