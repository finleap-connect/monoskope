package messaging

import (
	"context"
	"errors"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

// ErrCouldNotMarshalEvent is when an event could not be marshaled.
var ErrCouldNotMarshalEvent = errors.New("could not marshal event")

// ErrCouldNotPublishEvent is when cannot send event to message bus
var ErrCouldNotPublishEvent = errors.New("could not publish event")

// ErrMatcherMustNotBeNil is when an empty matcher has been provided
var ErrMatcherMustNotBeNil = errors.New("matcher must not be nil")

// ErrReceiverMustNotBeNil is when an empty receiver has been provided
var ErrReceiverMustNotBeNil = errors.New("receiver must not be nil")

// ErrMessageNotConnected is when there is no connection
var ErrMessageNotConnected = errors.New("message bus not connected")

// ErrMessageBusConnection is when an unexpected error on message bus occured
var ErrMessageBusConnection = errors.New("unexpected error on message bus occured")

// ErrCouldNotAddReceiver is when an receiver could not be added
var ErrCouldNotAddReceiver = errors.New("could not add receiver")

// ErrContextDeadlineExceeded is when execution has been aborted since the context deadline has been exceeded
var ErrContextDeadlineExceeded = errors.New("context deadline exceeded")

// EventBusPublisher publishes events on the underlying message bus.
type EventBusPublisher interface {
	// Connect connects to the bus
	Connect(context.Context) *messageBusError
	// PublishEvent publishes the event on the bus.
	PublishEvent(context.Context, events.Event) *messageBusError
	// Close closes the underlying connections
	Close() error
}

// EventBusConsumer notifies registered receivers on incoming events on the underlying message bus.
type EventBusConsumer interface {
	// Connect connects to the bus
	Connect(context.Context) *messageBusError
	// Matcher returns a new implementation specific matcher.
	Matcher() EventMatcher
	// AddReceiver adds a receiver for events matching one of the given EventMatcher.
	AddReceiver(context.Context, EventReceiver, ...EventMatcher) *messageBusError
	// Close closes the underlying connections
	Close() error
}

// EventMatcher is an interface used to define what events should be consumed
type EventMatcher interface {
	// Any matches any event.
	Any() EventMatcher
	// MatchEventType matches a specific event type, nil events never match.
	MatchEventType(eventType events.EventType) EventMatcher
	// MatchAggregate matches a specific aggregate type, nil events never match.
	MatchAggregateType(aggregateType events.AggregateType) EventMatcher
}

// EventReceiver is the function to call by the consumer on incoming events
type EventReceiver func(events.Event) error

// messageBusError is an error from the bus
type messageBusError struct {
	// Err is the error.
	Err error
	// BaseErr is an optional underlying error, for example from the message bus driver.
	BaseErr error
}

// Error implements the Error method of the errors.Error interface.
func (e messageBusError) Error() string {
	errStr := e.Err.Error()
	if e.BaseErr != nil {
		errStr += ": " + e.BaseErr.Error()
	}
	return errStr
}

// Cause returns the cause of this error.
func (e messageBusError) Cause() error {
	return e.Err
}

// UnwraMessageBusError returns the given error as MessageBusError if it is one
func UnwraMessageBusError(err error) *messageBusError {
	if esErr, ok := err.(messageBusError); ok {
		return &esErr
	}
	return nil
}
