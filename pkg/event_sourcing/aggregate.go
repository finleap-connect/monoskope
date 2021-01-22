package event_sourcing

import (
	"time"

	"github.com/google/uuid"
)

// EventType is the type of an event, used as its unique identifier.
type AggregateType string

// String returns the string representation of an AggregateType.
func (t AggregateType) String() string {
	return string(t)
}

type Aggregate interface {
	// AggregateType is the type of the aggregate that the event can be applied to.
	AggregateType() AggregateType
	// AggregateID is the id of the aggregate that the event should be applied to.
	AggregateID() uuid.UUID
	// AggregateVersion is the version of the aggregate.
	AggregateVersion() uint64
}

type AggregateBase struct {
	id            uuid.UUID
	aggregateType AggregateType
	version       uint64
	events        []Event
}

// NewAggregateBase creates an aggregate.
func NewAggregateBase(t AggregateType, id uuid.UUID) *AggregateBase {
	return &AggregateBase{
		id:            id,
		aggregateType: t,
	}
}

// EntityID implements the EntityID method of the Entity and Aggregate interface.
func (a *AggregateBase) AggregateID() uuid.UUID {
	return a.id
}

// AggregateType implements the AggregateType method of the Aggregate interface.
func (a *AggregateBase) AggregateType() AggregateType {
	return a.aggregateType
}

// Version implements the Version method of the Aggregate interface.
func (a *AggregateBase) AggregateVersion() uint64 {
	return a.version
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(et EventType, data EventData) Event {
	a.version++
	newEvent := NewEventFromAggregate(
		et,
		data,
		time.Now().UTC(),
		a)
	a.events = append(a.events, newEvent)

	return newEvent
}
