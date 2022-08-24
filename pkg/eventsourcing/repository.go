// Copyright 2022 Monoskope Authors
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

package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// RegistryObserver is an interface which must be implemented to register as observer for a Registry.
type RepositoryObserver[T Projection] interface {
	// Notify is called by the repository when an projection has been updated
	Notify(context.Context, T)
}

// Repository is a repository for reading projections.
type Repository[T Projection] interface {
	// ById returns a projection for an ID.
	ById(context.Context, uuid.UUID) (T, error)

	// All returns all projections in the repository.
	All(context.Context) ([]T, error)

	// Upsert saves a projection in the storage or replaces an existing one.
	Upsert(context.Context, T) error

	// RegisterObserver registers the given observer with the registry
	RegisterObserver(RepositoryObserver[T])

	// DeregisterObserver unregisters the given observer with the registry
	DeregisterObserver(RepositoryObserver[T])
}
