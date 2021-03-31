package eventsourcing

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type EventStreamReceiver interface {
	Receive() (Event, error)
}

type EventStreamSender interface {
	Send(Event)
	Error(error)
	Done()
}

type EventStream interface {
	EventStreamSender
	EventStreamReceiver
}
type eventStreamResult struct {
	events chan Event
	errors chan error
	done   chan int
}

func NewEventStream() EventStream {
	return &eventStreamResult{
		events: make(chan Event),
		errors: make(chan error),
		done:   make(chan int),
	}
}

func (e *eventStreamResult) Send(event Event) {
	if event != nil {
		e.events <- event
	}
}

func (e *eventStreamResult) Error(err error) {
	if err != nil {
		e.errors <- err
	}
}

func (e *eventStreamResult) Receive() (Event, error) {
	select {
	case <-e.done:
		return nil, io.EOF
	case event := <-e.events:
		return event, nil
	case err := <-e.errors:
		return nil, err
	}
}

func (e *eventStreamResult) Done() {
	defer close(e.done)
	defer close(e.events)
	defer close(e.errors)
	e.done <- 1
}

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
