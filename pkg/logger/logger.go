// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/operation"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

type Logger = logr.Logger

var (
	zapLog        *zap.Logger
	operationMode operation.OperationMode
	logMode       string
)

const (
	DebugLevel = 1
	InfoLevel  = 0
	WarnLevel  = -1
	ErrorLevel = -2
)

func init() {
	var err error

	// from build flag
	if logMode == "noop" {
		zapLog = zap.NewNop()
		return
	}

	// from env
	operationMode = operation.GetOperationMode()
	switch operationMode {
	case operation.DEVELOPMENT:
		zapLog, err = zap.NewDevelopment()
	case operation.RELEASE:
		zapLog, err = zap.NewProduction()
	default:
		zapLog = zap.NewNop()
	}

	if err != nil {
		panic(fmt.Sprintf("failed to setup logging: %v", err))
	}
}

func WithOptions(opts ...zap.Option) logr.Logger {
	switch operationMode {
	case operation.DEVELOPMENT:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(DebugLevel)
	case operation.RELEASE:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(InfoLevel)
	default:
		return zapr.NewLogger(zapLog.WithOptions(opts...)).V(ErrorLevel)
	}
}

func WithName(name string) logr.Logger {
	return WithOptions(zap.AddCaller()).WithName(name)
}
