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

	// LoadById loads the event by it's id from the store.
	LoadById(context.Context, uuid.UUID) (Event, error)
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

// EventData is any additional data for an event.
type EventData interface{}

// EventType is the type of an event, used as its unique identifier.
type EventType string

// EventType is the type of an event, used as its unique identifier.
type AggregateType string

// Event describes anything that has happened in the system.
// An event type name should be in past tense and contain the intent
// (TenantUpdated). The event should contain all the data needed when
// applying/handling it.
// The combination of aggregate_type, aggregate_id and version is
// unique.
type Event interface {
	// Type of the event.
	EventType() EventType
	// Type of the aggregate that the event can be applied to.
	AggregateType() AggregateType
	// ID of the aggregate that the event should be applied to.
	AggregateID() uuid.UUID
	// Timestamp of when the event was created.
	Timestamp() time.Time
	// Strict monotone counter, per aggregate/aggregate_id relation.
	Version() uint64
	// Event type specific event data.
	Data() EventData
	// A string representation of the event.
	String() string
}
