package eventsourcing

import (
	"context"
)

// Reactor is the interface for reactors.
type Reactor interface {
	// HandleEvent handles a given event returns 0..* Events in reaction or an error
	HandleEvent(context.Context, Event) ([]Event, error)
}
