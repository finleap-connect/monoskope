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
	log    Logger
	zapLog *zap.Logger
)

func init() {
	logMode := os.Getenv("LOG_MODE")
	var (
		err error
	)
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
	log = zapr.NewLogger(zapLog)
}

func Default() logr.Logger {
	return log
}

func WithOptions(opts ...zap.Option) logr.Logger {
	return zapr.NewLogger(zapLog.WithOptions(opts...))
}

func WithName(name string) logr.Logger {
	return WithOptions(zap.AddCaller()).WithName(name)
}
