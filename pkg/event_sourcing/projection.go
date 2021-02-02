package event_sourcing

import (
	"context"

	"github.com/google/uuid"
)

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the projection.
	ID() uuid.UUID
}

// Projector is the interface for projectors.
type Projector interface {
	// EvenTypes returns the EvenTypes for which events should be projected.
	EvenTypes() []EventType

	// AggregateTypes returns the AggregateTypes for which events should be projected.
	AggregateTypes() []AggregateType

	// Project updates the state of the projection occording to the given event.
	Project(context.Context, Event, Projection) (Projection, error)
}
