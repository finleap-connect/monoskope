package repositories

import (
	"context"
	"sync"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

// inMemoryRepository is a repository which stores projections in memory.
type inMemoryRepository struct {
	store map[string]es.Projection
	mutex sync.RWMutex
}

// NewInMemoryRepository creates a new repository which stores projections in memory.
func NewInMemoryRepository() es.ProjectionRepository {
	return &inMemoryRepository{
		store: make(map[string]es.Projection),
	}
}

// ById returns a projection for an ID.
func (r *inMemoryRepository) ById(ctx context.Context, id string) (es.Projection, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if val, ok := r.store[id]; ok {
		return val, nil
	}
	return nil, errors.ErrProjectionNotFound
}

// All returns all projections in the repository.
func (r *inMemoryRepository) All(context.Context) ([]es.Projection, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	all := make([]es.Projection, 0)
	for _, v := range r.store {
		all = append(all, v)
	}
	return all, nil
}

// Upsert saves a projection in the storage or replaces an existing one.
func (r *inMemoryRepository) Upsert(ctx context.Context, p es.Projection) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.store[p.GetId()] = p
	return nil
}

// Remove removes a projection by ID from the storage.
func (r *inMemoryRepository) Remove(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.store, id)
	return nil
}
