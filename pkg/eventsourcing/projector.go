package eventsourcing

import (
	"context"
)

// Projector is the interface for projectors.
type Projector interface {
	// NewProjection creates a new Projection of the type the Projector projects.
	NewProjection() Projection

	// Project updates the state of the projection occording to the given event.
	Project(context.Context, Event, Projection) (Projection, error)
}
