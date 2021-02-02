package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type inMemoryRepository struct {
	store map[uuid.UUID]es.Projection
}

func NewInMemoryRepository() es.Repository {
	return &inMemoryRepository{
		store: make(map[uuid.UUID]es.Projection),
	}
}

// Find returns an entity for an ID.
func (r *inMemoryRepository) Find(ctx context.Context, id uuid.UUID) (es.Projection, error) {
	if val, ok := r.store[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("not found")
}

// FindAll returns all entities in the repository.
func (r *inMemoryRepository) FindAll(context.Context) ([]es.Projection, error) {
	var all []es.Projection
	for _, v := range r.store {
		all = append(all, v)
	}
	return all, nil
}

// Save saves a entity in the storage.
func (r *inMemoryRepository) Save(ctx context.Context, p es.Projection) error {
	r.store[p.ID()] = p
	return nil
}

// Remove removes a entity by ID from the storage.
func (r *inMemoryRepository) Remove(ctx context.Context, id uuid.UUID) error {
	delete(r.store, id)
	return nil
}
