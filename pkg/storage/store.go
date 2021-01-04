package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNoEventsToAppend is when no events are available to append.
var ErrNoEventsToAppend = errors.New("no events to append")

// ErrIncorrectEventAggregateVersion is when an event is for an other version of the aggregate.
var ErrIncorrectAggregateVersion = errors.New("mismatching event aggreagte version")

// ErrAggregateVersionAlreadyExists is when an event is referencing an older version of the aggregate than is stored.
var ErrAggregateVersionAlreadyExists = errors.New("event aggreagte version already exists in store")

// ErrInvalidAggregateType is when an event is for a different type of aggregate.
var ErrInvalidAggregateType = errors.New("mismatching event aggreagte type")

// ErrCouldNotMarshalEvent is when an event could not be marshaled into JSON.
var ErrCouldNotMarshalEvent = errors.New("could not marshal event")

// ErrCouldNotMarshalEventContext is when an event could not be marshaled into JSON.
var ErrCouldNotMarshalEventContext = errors.New("could not marshal event context")

// ErrCouldNotSaveEvents is when events could not be saved.
var ErrCouldNotSaveEvents = errors.New("could not save events")

// Store is an interface for an event storage backend.
type Store interface {
	// Save appends all events in the event stream to the store.
	Save(context.Context, []Event) error

	// Load loads all events for the query from the store.
	Load(context.Context, *StoreQuery) ([]Event, error)

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

// EventStoreError is an error in the event store.
type EventStoreError struct {
	// Err is the error.
	Err error
	// BaseErr is an optional underlying error, for example from the DB driver.
	BaseErr error
}

// Error implements the Error method of the errors.Error interface.
func (e EventStoreError) Error() string {
	errStr := e.Err.Error()
	if e.BaseErr != nil {
		errStr += ": " + e.BaseErr.Error()
	}
	return errStr
}

// Cause returns the cause of this error.
func (e EventStoreError) Cause() error {
	return e.Err
}

// UnwrapEventStoreError returns the given error as EventStoreError if it is one
func UnwrapEventStoreError(err error) *EventStoreError {
	if esErr, ok := err.(EventStoreError); ok {
		return &esErr
	}
	return nil
}
