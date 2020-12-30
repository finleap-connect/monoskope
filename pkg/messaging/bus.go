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

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEvent matches a specific event type, nil events never match.
	MatchEvent(eventType storage.EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type, nil events never match.
	MatchAggregate(aggregateType storage.AggregateType) EventMatcher
}

// EventReceiver is the function to call by the consumer on incoming events
type EventReceiver func(storage.Event) error
