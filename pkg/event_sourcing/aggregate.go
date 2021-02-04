package event_sourcing

import (
	"time"

	"github.com/google/uuid"
)

// AggregateType is the type of an aggregate, used as its unique identifier.
type AggregateType string

// String returns the string representation of an AggregateType.
func (t AggregateType) String() string {
	return string(t)
}

// Aggregate is the interface definition for all aggregates
type Aggregate interface {
	// Type is the type of the aggregate that the event can be applied to.
	Type() AggregateType
	// ID is the id of the aggregate that the event should be applied to.
	ID() uuid.UUID
	// Version is the version of the aggregate.
	Version() uint64
	// Events returns the events that built up the aggregate.
	Events() []Event
}

// BaseAggregate is the base implementation for all aggregates
type BaseAggregate struct {
	id            uuid.UUID
	aggregateType AggregateType
	version       uint64
	events        []Event
}

// NewBaseAggregate creates an aggregate.
func NewBaseAggregate(t AggregateType, id uuid.UUID) *BaseAggregate {
	return &BaseAggregate{
		id:            id,
		aggregateType: t,
	}
}

// ID implements the ID method of the Aggregate interface.
func (a *BaseAggregate) ID() uuid.UUID {
	return a.id
}

// Type implements the Type method of the Aggregate interface.
func (a *BaseAggregate) Type() AggregateType {
	return a.aggregateType
}

// Version implements the Version method of the Aggregate interface.
func (a *BaseAggregate) Version() uint64 {
	return a.version
}

// Events implements the Events method of the Aggregate interface.
func (a *BaseAggregate) Events() []Event {
	return a.events
}

// AppendEvent appends an event to the events the aggregate was build upon.
func (a *BaseAggregate) AppendEvent(eventType EventType, eventData EventData) Event {
	a.version++
	newEvent := NewEventFromAggregate(
		eventType,
		eventData,
		time.Now().UTC(),
		a)
	a.events = append(a.events, newEvent)
	return newEvent
}
