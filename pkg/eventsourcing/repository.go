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
}

// Repository is a repository for reading and writing projections.
type Repository interface {
	ReadOnlyRepository
	WriteOnlyRepository
}
