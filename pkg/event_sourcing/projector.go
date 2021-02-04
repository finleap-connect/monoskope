package event_sourcing

import (
	"context"
)

// Projector is the interface for projectors.
type Projector interface {
	// AggregateType returns the AggregateType for which events should be projected.
	AggregateType() AggregateType

	// NewProjection creates a new Projection of the type the Projector projects.
	NewProjection() Projection

	// Project updates the state of the projection occording to the given event.
	Project(context.Context, Event, Projection) (Projection, error)
}
