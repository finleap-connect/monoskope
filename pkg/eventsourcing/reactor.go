package eventsourcing

import (
	"context"
)

// Reactor is the interface for reactors.
type Reactor interface {
	// HandleEvent handles a given event send 0..* Events through the given channel in reaction or an error.
	// Attention: The reactor is responsible for closing the channel if no further events will be send to that channel.
	HandleEvent(ctx context.Context, event Event, events chan<- Event) error
}
