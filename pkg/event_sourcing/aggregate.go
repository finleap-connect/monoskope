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
	// AggregateType returns the type of the aggregate.
	AggregateType() AggregateType
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

// EntityID implements the EntityID method of the eh.Entity and eh.Aggregate interface.
func (a *AggregateBase) EntityID() uuid.UUID {
	return a.id
}

// AggregateType implements the AggregateType method of the eh.Aggregate interface.
func (a *AggregateBase) AggregateType() AggregateType {
	return a.aggregateType
}

// Version implements the Version method of the Aggregate interface.
func (a *AggregateBase) Version() uint64 {
	return a.version
}

// IncrementVersion implements the IncrementVersion method of the Aggregate interface.
func (a *AggregateBase) IncrementVersion() {
	a.version++
}

// Events implements the Events method of the eh.EventSource interface.
func (a *AggregateBase) Events() []Event {
	events := a.events
	a.events = nil
	return events
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(et EventType, data EventData) Event {
	newEvent := NewEvent(
		et,
		data,
		time.Now().UTC(),
		a.AggregateType(),
		a.EntityID(),
		a.Version()+uint64(len(a.events)+1))

	a.events = append(a.events, newEvent)
	return newEvent
}
