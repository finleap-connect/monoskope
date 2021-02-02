package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// inMemoryRepository is a repository which stores projections in memory.
type inMemoryRepository struct {
	store map[uuid.UUID]es.Projection
}

// NewInMemoryRepository creates a new repository which stores projections in memory.
func NewInMemoryRepository() es.Repository {
	return &inMemoryRepository{
		store: make(map[uuid.UUID]es.Projection),
	}
}

// ById returns a projection for an ID.
func (r *inMemoryRepository) ById(ctx context.Context, id uuid.UUID) (es.Projection, error) {
	if val, ok := r.store[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("not found")
}

// All returns all projections in the repository.
func (r *inMemoryRepository) All(context.Context) ([]es.Projection, error) {
	all := make([]es.Projection, 0)
	for _, v := range r.store {
		all = append(all, v)
	}
	return all, nil
}

// Upsert saves a projection in the storage or replaces an existing one.
func (r *inMemoryRepository) Upsert(ctx context.Context, p es.Projection) error {
	r.store[p.ID()] = p
	return nil
}

// Remove removes a projection by ID from the storage.
func (r *inMemoryRepository) Remove(ctx context.Context, id uuid.UUID) error {
	delete(r.store, id)
	return nil
}
