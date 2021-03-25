package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// AggregateStore handles storing and loading Aggregates.
type AggregateManager interface {
	// Get returns the most recent version of all aggregate of a given type.
	All(context.Context, AggregateType, bool) ([]Aggregate, error)

	// Get returns the most recent version of an aggregate.
	Get(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Update stores all in-flight events for an aggregate.
	Update(context.Context, Aggregate) error
}
