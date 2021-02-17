package eventsourcing

import "github.com/google/uuid"

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the Projection.
	ID() uuid.UUID
	// Version returns the version of the aggregate this Projection is based upon.
	Version() uint64
	// IncrementVersion increments the Version of the Projection.
	IncrementVersion()
}
