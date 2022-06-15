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

package repositories

import (
	"context"
	"sync"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

// inMemoryRepository is a repository which stores projections in memory.
type inMemoryRepository[T es.Projection] struct {
	store map[uuid.UUID]T
	mutex sync.RWMutex
}

// NewInMemoryRepository creates a new repository which stores projections in memory.
func NewInMemoryRepository[T es.Projection]() es.Repository[T] {
	return &inMemoryRepository[T]{
		store: make(map[uuid.UUID]T),
	}
}

// ById returns a projection for an ID.
func (r *inMemoryRepository[T]) ById(ctx context.Context, id uuid.UUID) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.store[id]; ok {
		return val, nil
	}
	var result T
	return result, errors.ErrProjectionNotFound
}

// All returns all projections in the repository.
func (r *inMemoryRepository[T]) All(context.Context) ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	all := make([]T, 0)
	for _, v := range r.store {
		all = append(all, v)
	}
	return all, nil
}

// Upsert saves a projection in the storage or replaces an existing one.
func (r *inMemoryRepository[T]) Upsert(ctx context.Context, p T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.store[p.ID()] = p
	return nil
}

// Remove removes a projection by ID from the storage.
func (r *inMemoryRepository[T]) Remove(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.store, id)
	return nil
}
