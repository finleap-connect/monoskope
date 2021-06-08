package eventsourcing

import (
	"context"
)

// Reactor is the interface for reactors.
type Reactor interface {
	// HandleEvent handles a given event send 0..* Events through the given channel in reaction or an error
	HandleEvent(ctx context.Context, event Event, events chan<- Event) error
}
