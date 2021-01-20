package event_sourcing

import (
	"context"
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
	// AggregateType returns the type name of the aggregate.
	// AggregateType() string
	AggregateType() AggregateType

	// HandleCommand implements the HandleCommand method of the Aggregate.
	HandleCommand(context.Context, Command) error
}

type AggregateBase struct {
	id     uuid.UUID
	t      AggregateType
	v      uint64
	events []Event
}

// EntityID implements the EntityID method of the eh.Entity and eh.Aggregate interface.
func (a *AggregateBase) EntityID() uuid.UUID {
	return a.id
}

// AggregateType implements the AggregateType method of the eh.Aggregate interface.
func (a *AggregateBase) AggregateType() AggregateType {
	return a.t
}

// Version implements the Version method of the Aggregate interface.
func (a *AggregateBase) Version() uint64 {
	return a.v
}

// IncrementVersion implements the IncrementVersion method of the Aggregate interface.
func (a *AggregateBase) IncrementVersion() {
	a.v++
}

// Events implements the Events method of the eh.EventSource interface.
func (a *AggregateBase) Events() []Event {
	events := a.events
	a.events = nil
	return events
}

// AppendEvent appends an event for later retrieval by Events().
func (a *AggregateBase) AppendEvent(t EventType, data EventData, timestamp time.Time) Event {
	ev := NewEvent(
		t,
		data,
		timestamp,
		a.AggregateType(),
		a.EntityID(),
		a.Version()+uint64(len(a.events)+1))

	a.events = append(a.events, ev)
	return ev
}
