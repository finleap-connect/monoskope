package event_sourcing

import (
	"github.com/google/uuid"
)

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the Projection.
	ID() uuid.UUID
	// AggregateVersion is the version of the Aggregate this Projection is based upon.
	AggregateVersion() uint64
}

type BaseProjection struct {
	id      uuid.UUID
	version uint64
}

func NewBaseProjection(id uuid.UUID) BaseProjection {
	return BaseProjection{
		id: id,
	}
}

func (p BaseProjection) ID() uuid.UUID {
	return p.id
}

func (p BaseProjection) AggregateVersion() uint64 {
	return p.version
}
