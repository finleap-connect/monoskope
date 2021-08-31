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
