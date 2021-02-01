package event_sourcing

import (
	"context"

	"github.com/google/uuid"
)

// ReadOnlyRepository is a repository for reading projections.
type ReadOnlyRepository interface {
	// Find returns an entity for an ID.
	Find(context.Context, uuid.UUID) (Projection, error)

	// FindAll returns all entities in the repository.
	FindAll(context.Context) ([]Projection, error)
}

// WriteOnlyRepository is a repository for writing projections.
type WriteOnlyRepository interface {
	// Save saves a entity in the storage.
	Save(context.Context, Projection) error

	// Remove removes a entity by ID from the storage.
	Remove(context.Context, uuid.UUID) error
}

// Repository is a repository for reading and writing projections.
type Repository interface {
	ReadOnlyRepository
	WriteOnlyRepository
}
