package eventsourcing

import (
	"context"
)

// EventBusConnector can open and close connections.
type EventBusConnector interface {
	// Open connects to the bus
	Open(context.Context) error
	// Close closes the underlying connections
	Close() error
}

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	EventBusConnector
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, Event) error
}

// EventBusConsumer notifies registered handlers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	EventBusConnector
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddHandler adds a handler for events matching one of the given EventMatcher.
	AddHandler(context.Context, EventHandler, ...EventMatcher) error
}

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEventType matches a specific event type.
	MatchEventType(eventType EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type.
	MatchAggregateType(aggregateType AggregateType) EventMatcher
}
