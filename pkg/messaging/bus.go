package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// EventBus publishes events on the underlying message bus or notifies receivers on incoming events.
type EventBus interface {
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, storage.Event) error

	// AddReceiver adds a receiver for event matching the EventFilter.
	AddReceiver(EventMatcher, EventReceiver)
}
