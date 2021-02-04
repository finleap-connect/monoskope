package event_sourcing

import (
	"context"
)

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	// Connect connects to the bus
	Connect(context.Context) error
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, Event) error
	// Close closes the underlying connections
	Close() error
}

// EventBusConsumer notifies registered handlers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	// Connect connects to the bus
	Connect(context.Context) error
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddHandler adds a handler for events matching one of the given EventMatcher.
	AddHandler(context.Context, EventHandler, ...EventMatcher) error
	// Close closes the underlying connections
	Close() error
}

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEventType matches a specific event type, nil events never match.
	MatchEventType(eventType EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type, nil events never match.
	MatchAggregateType(aggregateType AggregateType) EventMatcher
}
