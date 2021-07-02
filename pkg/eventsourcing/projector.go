package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

// Projector is the interface for projectors.
type Projector interface {
	// NewProjection creates a new Projection of the type the Projector projects.
	NewProjection(uuid.UUID) Projection

	// Project updates the state of the projection according to the given event.
	Project(context.Context, Event, Projection) (Projection, error)
}
