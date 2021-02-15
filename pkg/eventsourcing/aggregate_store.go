package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// AggregateStore handles storing and loading Aggregates.
type AggregateHandler interface {
	// Get returns the most recent version of an aggregate.
	Get(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Update stores all in-flight events for an aggregate.
	Update(context.Context, Aggregate) error
}
