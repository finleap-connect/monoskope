package usecase

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// UseCase is the interface an UseCase must implement to be executable
type UseCase interface {
	// Run starts the execution of the behavior of a use case
	Run(context.Context) error
}

// UseCaseBase is the basic fields needed to implement a use case along with a logger
type UseCaseBase struct {
	Log logger.Logger
}

// NewUseCaseBase returns a new basic use case implementation
func NewUseCaseBase(name string) *UseCaseBase {
	return &UseCaseBase{
		Log: logger.WithName(fmt.Sprintf("%s-usecase", name)),
	}
}
