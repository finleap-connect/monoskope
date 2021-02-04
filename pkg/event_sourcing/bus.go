package event_sourcing

import (
	"context"
	"errors"
)

// ErrCouldNotMarshalEvent is when an event could not be marshaled.
var ErrCouldNotMarshalEvent = errors.New("could not marshal event")

// ErrCouldNotPublishEvent is when cannot send event to message bus
var ErrCouldNotPublishEvent = errors.New("could not publish event")

// ErrMatcherMustNotBeNil is when an empty matcher has been provided
var ErrMatcherMustNotBeNil = errors.New("matcher must not be nil")

// ErrHandlerMustNotBeNil is when an empty handler has been provided
var ErrHandlerMustNotBeNil = errors.New("handler must not be nil")

// ErrMessageNotConnected is when there is no connection
var ErrMessageNotConnected = errors.New("message bus not connected")

// ErrMessageBusConnection is when an unexpected error on message bus occured
var ErrMessageBusConnection = errors.New("unexpected error on message bus occured")

// ErrCouldNotAddHandler is when an handler could not be added
var ErrCouldNotAddHandler = errors.New("could not add handler")

// ErrContextDeadlineExceeded is when execution has been aborted since the context deadline has been exceeded
var ErrContextDeadlineExceeded = errors.New("context deadline exceeded")

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
