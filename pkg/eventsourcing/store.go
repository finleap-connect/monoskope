package eventsourcing

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

// Store is an interface for an event storage backend.
type Store interface {
	// Connect connects to the bus
	Connect(context.Context) error

	// Save appends all events in the event stream to the store.
	Save(context.Context, []Event) error

	// Load loads all events for the query from the store.
	Load(context.Context, *StoreQuery) (EventStreamReceiver, error)

	// Close closes the underlying connections
	Close() error
}

// StoreQuery contains query information on how to retrieve events from an event store
type StoreQuery struct {
	// Filter events by aggregate id
	AggregateId *uuid.UUID
	// Filter events for a specific aggregate type
	AggregateType *AggregateType
	// Filter events with a Version >= MinVersion
	MinVersion *uint64
	// Filter events with a Version <= MaxVersion
	MaxVersion *uint64
	// Filter events with a Timestamp >= MinTimestamp
	MinTimestamp *time.Time
	// Filter events with a Timestamp <= MaxTimestamp
	MaxTimestamp *time.Time
	// Filter events by aggregates that have not been deleted
	ExcludeDeleted bool
}

type EventStreamReceiver interface {
	Receive() (Event, error)
}

type EventStreamSender interface {
	Send(Event)
	Error(error)
	Done()
}

// EventStream is an interface for handling asynchronous result sending/receiving from the EventStore
type EventStream interface {
	EventStreamSender
	EventStreamReceiver
}

type eventStreamResult struct {
	events chan Event
	errors chan error
	done   chan int
}

// NewEventStream returns an implementation for the EventStream interface
func NewEventStream() EventStream {
	return &eventStreamResult{
		events: make(chan Event),
		errors: make(chan error),
		done:   make(chan int),
	}
}

// Send sends the given Event to the receiver
func (e *eventStreamResult) Send(event Event) {
	if event != nil {
		e.events <- event
	}
}

// Error sends the given error to the receiver
func (e *eventStreamResult) Error(err error) {
	if err != nil {
		e.errors <- err
	}
}

// Receive receives an Event or error from the sender.
// When the error io.EOF is received the stream has been closed by the sender.
func (e *eventStreamResult) Receive() (Event, error) {
	select {
	case <-e.done:
		return nil, io.EOF
	case err := <-e.errors:
		return nil, err
	case event := <-e.events:
		return event, nil
	}
}

// Done closes all channels
func (e *eventStreamResult) Done() {
	defer close(e.done)
	defer close(e.events)
	defer close(e.errors)
	e.done <- 1
}
