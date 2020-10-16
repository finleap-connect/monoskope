package logger

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

const (
	DebugLevel = 3
)

type Logger = logr.Logger

var log Logger

func init() {
	logMode := os.Getenv("LOG_MODE")
	var (
		zapLog *zap.Logger
		err    error
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

func WithName(name string) logr.Logger {
	return log.WithName(name)
}
