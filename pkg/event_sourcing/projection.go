package event_sourcing

import (
	"github.com/google/uuid"
)

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the projection.
	ID() uuid.UUID
}
