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

// Store is an interface for an event storage backend.
type Store interface {
	// Save appends all events in the event stream to the store.
	Save(context.Context, []Event) error

	// Load loads all events for the query from the store.
	Load(context.Context, *StoreQuery) ([]Event, error)
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

// EventStoreError is an error in the event store, with the namespace.
type EventStoreError struct {
	// Err is the error.
	Err error
	// BaseErr is an optional underlying error, for example from the DB driver.
	BaseErr error
}

func UnwrapEventStoreError(err error) *EventStoreError {
	if esErr, ok := err.(EventStoreError); ok {
		return &esErr
	}
	return nil
}
