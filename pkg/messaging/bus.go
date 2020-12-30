package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, storage.Event) error
}

// EventBusConsumer notifies registered receivers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddReceiver adds a receiver for event matching the EventFilter.
	AddReceiver(EventMatcher, EventReceiver) error
}
