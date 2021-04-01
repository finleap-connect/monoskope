package eventsourcing

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	// EventType is the type of the event.
	EventType() EventType
	// Timestamp of when the event was created.
	Timestamp() time.Time
	// AggregateType is the type of the aggregate that the event can be applied to.
	AggregateType() AggregateType
	// AggregateID is the id of the aggregate that the event should be applied to.
	AggregateID() uuid.UUID
	// AggregateVersion is the version of the aggregate.
	AggregateVersion() uint64
	// Event type specific event data.
	Data() EventData
	// Metadata is app-specific metadata originating user etc. when this event has been stored.
	Metadata() map[string]string
	// A string representation of the event.
	String() string
}

// NewEvent creates a new event.
func NewEvent(ctx context.Context, eventType EventType, data EventData, timestamp time.Time,
	aggregateType AggregateType, aggregateID uuid.UUID, aggregateVersion uint64) Event {

	return &event{
		eventType:        eventType,
		data:             data,
		timestamp:        timestamp,
		aggregateType:    aggregateType,
		aggregateID:      aggregateID,
		aggregateVersion: aggregateVersion,
		metadata:         NewMetadataManagerFromContext(ctx).GetMetadata(),
	}
}

// NewEvent creates a new event with metadata attached.
func NewEventWithMetadata(eventType EventType, data EventData, timestamp time.Time,
	aggregateType AggregateType, aggregateID uuid.UUID, aggregateVersion uint64, metadata map[string]string) Event {
	return event{
		eventType:        eventType,
		data:             data,
		timestamp:        timestamp,
		aggregateType:    aggregateType,
		aggregateID:      aggregateID,
		aggregateVersion: aggregateVersion,
		metadata:         metadata,
	}
}

// NewEventFromProto converts proto events to Event
func NewEventFromProto(protoEvent *esApi.Event) (Event, error) {
	aggregateId, err := uuid.Parse(protoEvent.GetAggregateId())
	if err != nil {
		return nil, errors.ErrCouldNotParseAggregateId
	}

	return NewEventWithMetadata(
		EventType(protoEvent.GetType()),
		protoEvent.GetData(),
		protoEvent.Timestamp.AsTime(),
		AggregateType(protoEvent.GetAggregateType()),
		aggregateId,
		protoEvent.GetAggregateVersion().GetValue(),
		protoEvent.Metadata,
	), nil
}

// NewProtoFromEvent converts Event to proto events
func NewProtoFromEvent(storeEvent Event) *esApi.Event {
	ev := &esApi.Event{
		Type:             storeEvent.EventType().String(),
		Timestamp:        timestamppb.New(storeEvent.Timestamp()),
		AggregateType:    storeEvent.AggregateType().String(),
		AggregateId:      storeEvent.AggregateID().String(),
		AggregateVersion: &wrapperspb.UInt64Value{Value: storeEvent.AggregateVersion()},
		Data:             storeEvent.Data(),
		Metadata:         storeEvent.Metadata(),
	}
	return ev
}

// event is an internal representation of an event.
type event struct {
	eventType        EventType
	data             EventData
	timestamp        time.Time
	aggregateType    AggregateType
	aggregateID      uuid.UUID
	aggregateVersion uint64
	metadata         map[string]string
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

// Metadata implements the Metadata method of the Event interface.
func (e event) Metadata() map[string]string {
	return e.metadata
}

// String implements the String method of the Event interface.
func (e event) String() string {
	return fmt.Sprintf("%s@%d", e.eventType, e.aggregateVersion)
}
