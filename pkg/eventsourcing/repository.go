package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// ReadOnlyProjectionRepository is a repository for reading projections.
type ReadOnlyProjectionRepository interface {
	// ById returns a projection for an ID.
	ById(context.Context, string) (Projection, error)

	// All returns all projections in the repository.
	All(context.Context) ([]Projection, error)
}

// WriteOnlyProjectionRepository is a repository for writing projections.
type WriteOnlyProjectionRepository interface {
	// Upsert saves a projection in the storage or replaces an existing one.
	Upsert(context.Context, Projection) error

	// Remove removes a projection by ID from the storage.
	Remove(context.Context, string) error
}

// ProjectionRepository is a repository for reading and writing projections.
type ProjectionRepository interface {
	ReadOnlyProjectionRepository
	WriteOnlyProjectionRepository
}

// AggregateStore is responsible for loading and saving aggregates.
type AggregateRepository interface {
	// Load loads the most recent version of an aggregate with a type and id.
	Load(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Save stores all unsaved events for an aggregate.
	Save(context.Context, Aggregate) error
}
