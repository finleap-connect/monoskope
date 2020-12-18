package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
