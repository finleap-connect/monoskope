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
type LogLevel = int

var (
	zapLog        *zap.Logger
	operationMode operation.OperationMode
	logMode       string
)

const (
	DebugLevel LogLevel = 1
	InfoLevel  LogLevel = 0
	WarnLevel  LogLevel = -1
	ErrorLevel LogLevel = -2
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
	return zapr.NewLogger(zapLog.WithOptions(opts...))
}

func WithName(name string) logr.Logger {
	return WithOptions(zap.AddCaller()).WithName(name)
}

type grpcLog struct {
	log   Logger
	level LogLevel
}

func (l *grpcLog) Write(p []byte) (n int, err error) {
	message := string(p)
	switch l.level {
	case InfoLevel:
		l.log.V(DebugLevel).WithValues("level", InfoLevel).Info(message)
	case WarnLevel:
		l.log.WithValues("level", WarnLevel).Info(message)
	case ErrorLevel:
		l.log.WithValues("level", ErrorLevel).Error(fmt.Errorf(message), message)
	}
	return len(p), nil
}

func NewGrpcLog(log Logger, level LogLevel) *grpcLog {
	return &grpcLog{log: log, level: level}
}
