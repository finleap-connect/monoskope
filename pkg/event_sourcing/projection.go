package event_sourcing

import (
	"github.com/google/uuid"
)

type Projection interface {
	// ID returns the ID of the projection.
	ID() uuid.UUID
}

type Projector interface {
	Matchers() []*EventMatcher
}
