package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/operation"
	"go.uber.org/zap"
)

type Logger = logr.Logger

var (
	zapLog *zap.Logger
)

func init() {
	var err error

	operationMode := operation.GetOperationMode()

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
