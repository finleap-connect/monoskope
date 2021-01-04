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

// ErrMessageBusConnection is when an unexpected error on message bus occured
var ErrMessageBusConnection = errors.New("unexpected error on message bus occured")

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, storage.Event) *MessageBusError
	// Close for freeing all disposable resources
	Close() error
}

// EventBusConsumer notifies registered receivers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddReceiver adds a receiver for events matching one of the given EventMatcher.
	AddReceiver(EventReceiver, ...EventMatcher) error
	// AddErrorHandler adds a handler function to call if any error occurs
	AddErrorHandler(ErrorHandler)
	// Close frees all disposable resources
	Close() error
}

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEventType matches a specific event type, nil events never match.
	MatchEventType(eventType storage.EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type, nil events never match.
	MatchAggregateType(aggregateType storage.AggregateType) EventMatcher
}

// EventReceiver is the function to call by the consumer on incoming events
type EventReceiver func(storage.Event) error

type ErrorHandler func(MessageBusError)

// MessageBusError is an error from the bus
type MessageBusError struct {
	// Err is the error.
	Err error
	// BaseErr is an optional underlying error, for example from the message bus driver.
	BaseErr error
}

// Error implements the Error method of the errors.Error interface.
func (e MessageBusError) Error() string {
	errStr := e.Err.Error()
	if e.BaseErr != nil {
		errStr += ": " + e.BaseErr.Error()
	}
	return errStr
}

// Cause returns the cause of this error.
func (e MessageBusError) Cause() error {
	return e.Err
}

// UnwraMessageBusError returns the given error as MessageBusError if it is one
func UnwraMessageBusError(err error) *MessageBusError {
	if esErr, ok := err.(MessageBusError); ok {
		return &esErr
	}
	return nil
}
