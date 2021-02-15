package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// ReadOnlyProjectionRepository is a repository for reading projections.
type ReadOnlyProjectionRepository interface {
	// ById returns a projection for an ID.
	ById(context.Context, uuid.UUID) (Projection, error)

	// All returns all projections in the repository.
	All(context.Context) ([]Projection, error)
}

// WriteOnlyProjectionRepository is a repository for writing projections.
type WriteOnlyProjectionRepository interface {
	// Upsert saves a projection in the storage or replaces an existing one.
	Upsert(context.Context, Projection) error

	// Remove removes a projection by ID from the storage.
	Remove(context.Context, uuid.UUID) error
}

// ProjectionRepository is a repository for reading and writing projections.
type ProjectionRepository interface {
	ReadOnlyProjectionRepository
	WriteOnlyProjectionRepository
}

// AggregateRepository is a repository for reading and writing aggregates.
type AggregateRepository interface {
	// Get returns the most recent version of an aggregate.
	Get(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Update stores all in-flight events for an aggregate.
	Update(context.Context, Aggregate) error
}
