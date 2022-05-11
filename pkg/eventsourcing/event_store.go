// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventsourcing

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

// EventStore is an interface for an event storage backend.
type EventStore interface {
	// Open connects to the bus
	Open(context.Context) error

	// Save appends all events in the event stream to the store.
	Save(context.Context, []Event) error

	// Load loads all events for the query from the store.
	Load(context.Context, *StoreQuery) (EventStreamReceiver, error)

	// LoadOr loads all events by combining the queries with the logical OR from the store.
	LoadOr(context.Context, []*StoreQuery) (EventStreamReceiver, error)

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
