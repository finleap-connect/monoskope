package messaging

import (
	"context"
	"errors"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// ErrCouldNotMarshalEvent is when an event could not be marshaled.
var ErrCouldNotMarshalEvent = errors.New("could not marshal event")

// ErrCouldNotUnmarshalEvent is when an event could not be unmarshaled.
var ErrCouldNotUnmarshalEvent = errors.New("could not unmarshal event")

// ErrCouldNotPublishEvent is when cannot send event to message bus
var ErrCouldNotPublishEvent = errors.New("could not publish event")

// ErrMatcherMustNotBeNil is when an empty matcher has been provided
var ErrMatcherMustNotBeNil = errors.New("matcher must not be nil")

// ErrReceiverMustNotBeNil is when an empty receiver has been provided
var ErrReceiverMustNotBeNil = errors.New("receiver must not be nil")

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
