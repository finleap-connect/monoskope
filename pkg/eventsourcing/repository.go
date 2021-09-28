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

package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// ReadOnlyRepository is a repository for reading projections.
type ReadOnlyRepository interface {
	// ById returns a projection for an ID.
	ById(context.Context, uuid.UUID) (Projection, error)

	// All returns all projections in the repository.
	All(context.Context) ([]Projection, error)
}

// WriteOnlyRepository is a repository for writing projections.
type WriteOnlyRepository interface {
	// Upsert saves a projection in the storage or replaces an existing one.
	Upsert(context.Context, Projection) error

	// Remove removes a projection by ID from the storage.
	Remove(context.Context, uuid.UUID) error
}

// Repository is a repository for reading and writing projections.
type Repository interface {
	ReadOnlyRepository
	WriteOnlyRepository
}
