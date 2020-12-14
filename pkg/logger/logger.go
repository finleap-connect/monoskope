package logger

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

type Logger = logr.Logger

var (
	zapLog  *zap.Logger
	logMode string
)

func init() {
	var err error

	if logMode == "" {
		logMode = os.Getenv("LOG_MODE")
	}

	if logMode == "" || logMode == "dev" {
		zapLog, err = zap.NewDevelopment()
	} else if logMode == "prod" {
		zapLog, err = zap.NewProduction()
	} else {
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
