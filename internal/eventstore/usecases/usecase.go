package usecases

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// UseCase is the interface an UseCase must implement to be executable
type UseCase interface {
	// Run starts the execution of the behaviour of a use case
	Run() error
}

// UseCaseBase is the basic fields needed to implement a use case along with a logger and a context
type UseCaseBase struct {
	log logger.Logger
	ctx context.Context
}
