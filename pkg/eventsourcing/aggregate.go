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
	CommandHandler
	setId(id uuid.UUID)
	// Type is the type of the aggregate that the event can be applied to.
	Type() AggregateType
	// ID is the id of the aggregate that the event should be applied to.
	ID() uuid.UUID
	// Version is the version of the aggregate.
	Version() uint64
	// SetDeleted sets the deleted flag of the aggregate to true
	SetDeleted(bool)
	// Deleted indicates whether the aggregate resource has been deleted
	Deleted() bool
	// UncommittedEvents returns outstanding events that need to persisted. They are cleared on reading them.
	UncommittedEvents() []Event
	// ApplyEvent applies an Event on the aggregate.
	ApplyEvent(Event) error
	// IncrementVersion increments the version of the Aggregate.
	IncrementVersion()
	// Exists returns if the version of the aggregate is >0
	Exists() bool
}

// BaseAggregate is the base implementation for all aggregates
type BaseAggregate struct {
	id            uuid.UUID
	aggregateType AggregateType
	version       uint64
	deleted       bool
	events        []Event
}

// NewBaseAggregate creates an aggregate.
func NewBaseAggregate(t AggregateType) *BaseAggregate {
	return &BaseAggregate{
		id:            uuid.New(),
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

// SetDeleted implements the SetDeleted method of the Aggregate interface.
func (a *BaseAggregate) SetDeleted(deleted bool) {
	a.deleted = deleted
}

// Deleted implements the Deleted method of the Aggregate interface.
func (a *BaseAggregate) Deleted() bool {
	return a.deleted
}

// Exists returns if the version of the aggregate is >0
func (a *BaseAggregate) Exists() bool {
	return a.version > 0
}

// UncommittedEvents implements the UncommittedEvents method of the Aggregate interface.
func (a *BaseAggregate) UncommittedEvents() []Event {
	defer func() {
		a.events = nil
	}()
	return a.events
}

// AppendEvent appends an event to the events the aggregate was build upon.
func (a *BaseAggregate) AppendEvent(ctx context.Context, eventType EventType, eventData EventData) Event {
	newEvent := NewEvent(
		ctx,
		eventType,
		eventData,
		time.Now().UTC(),
		a.Type(),
		a.ID(),
		a.version+uint64(len(a.events))+1)
	a.events = append(a.events, newEvent)
	return newEvent
}

// IncrementVersion implements the IncrementVersion method of the Aggregate interface.
func (a *BaseAggregate) IncrementVersion() {
	a.version++
}

// setId implements the private method to set Aggregate id.
func (a *BaseAggregate) setId(id uuid.UUID) {
	a.id = id
}
