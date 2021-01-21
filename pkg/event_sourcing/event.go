package event_sourcing

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EventType is the type of an event, used as its unique identifier.
type EventType string

// String returns the string representation of an EventType.
func (t EventType) String() string {
	return string(t)
}

// Event describes anything that has happened in the system.
// An event type name should be in past tense and contain the intent
// (TenantUpdated). The event should contain all the data needed when
// applying/handling it.
// The combination of AggregateType, AggregateID and AggregateVersion is
// unique.
type Event interface {
	// Type of the event.
	EventType() EventType
	// Type of the aggregate that the event can be applied to.
	Timestamp() time.Time
	// Strict monotone counter, per aggregate/aggregate_id relation.
	AggregateType() AggregateType
	// ID of the aggregate that the event should be applied to.
	AggregateID() uuid.UUID
	// Timestamp of when the event was created.
	AggregateVersion() uint64
	// Event type specific event data.
	Data() EventData
	// A string representation of the event.
	String() string
}

// NewEvent creates a new event with a type and data, setting its timestamp.
func NewEvent(eventType EventType, data EventData, timestamp time.Time,
	aggregateType AggregateType, aggregateID uuid.UUID, aggregateVersion uint64) Event {
	return event{
		eventType:        eventType,
		data:             data,
		timestamp:        timestamp,
		aggregateType:    aggregateType,
		aggregateID:      aggregateID,
		aggregateVersion: aggregateVersion,
	}
}

// event is an internal representation of an event, returned when the aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type event struct {
	eventType        EventType
	data             EventData
	timestamp        time.Time
	aggregateType    AggregateType
	aggregateID      uuid.UUID
	aggregateVersion uint64
}

// EventType implements the EventType method of the Event interface.
func (e event) EventType() EventType {
	return e.eventType
}

// Data implements the Data method of the Event interface.
func (e event) Data() EventData {
	return e.data
}

// Timestamp implements the Timestamp method of the Event interface.
func (e event) Timestamp() time.Time {
	return e.timestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e event) AggregateType() AggregateType {
	return e.aggregateType
}

// AggrgateID implements the AggrgateID method of the Event interface.
func (e event) AggregateID() uuid.UUID {
	return e.aggregateID
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e event) AggregateVersion() uint64 {
	return e.aggregateVersion
}

// String implements the String method of the Event interface.
func (e event) String() string {
	return fmt.Sprintf("%s@%d", e.eventType, e.aggregateVersion)
}
