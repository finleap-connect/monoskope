package usecases

import (
	"errors"

	"github.com/google/uuid"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ErrCouldNotParseAggregateId is when an aggregate id could not be parsed as uuid
var ErrCouldNotParseAggregateId = errors.New("could not parse aggregate id")

// NewEventFromProto converts api_es.Event to evs.Event
func NewEventFromProto(protoEvent *api_es.Event) (evs.Event, error) {
	aggregateId, err := uuid.Parse(protoEvent.GetAggregateId())
	if err != nil {
		return nil, ErrCouldNotParseAggregateId
	}

	ev := evs.NewEvent(
		evs.EventType(protoEvent.GetType()),
		protoEvent.GetData(),
		protoEvent.Timestamp.AsTime(),
		evs.AggregateType(protoEvent.GetAggregateType()),
		aggregateId,
		protoEvent.GetAggregateVersion().GetValue())
	return ev, nil
}

// NewStoreQueryFromProto converts proto api_es.EventFilter to storage.StoreQuery
func NewStoreQueryFromProto(protoFilter *api_es.EventFilter) (*storage.StoreQuery, error) {
	storeQuery := &storage.StoreQuery{}

	if val, ok := protoFilter.GetByAggregate().(*api_es.EventFilter_AggregateId); ok {
		aId, err := uuid.Parse(val.AggregateId.GetValue())
		if err != nil {
			return nil, ErrCouldNotParseAggregateId
		}
		storeQuery.AggregateId = &aId
	}
	if val, ok := protoFilter.GetByAggregate().(*api_es.EventFilter_AggregateType); ok {
		aType := evs.AggregateType(val.AggregateType.GetValue())
		storeQuery.AggregateType = &aType
	}

	if protoFilter.GetMinVersion() != nil {
		val := protoFilter.GetMinVersion().GetValue()
		storeQuery.MinVersion = &val
	}
	if protoFilter.GetMaxVersion() != nil {
		val := protoFilter.GetMaxVersion().GetValue()
		storeQuery.MaxVersion = &val
	}

	if protoFilter.GetMinTimestamp() != nil {
		val := protoFilter.GetMinTimestamp().AsTime()
		storeQuery.MinTimestamp = &val
	}
	if protoFilter.GetMaxTimestamp() != nil {
		val := protoFilter.GetMaxTimestamp().AsTime()
		storeQuery.MaxTimestamp = &val
	}

	return storeQuery, nil
}

// NewProtoFromEvent converts evs.Event to api_es.Event
func NewProtoFromEvent(storeEvent evs.Event) (*api_es.Event, error) {
	ev := &api_es.Event{
		Type:             storeEvent.EventType().String(),
		Timestamp:        timestamppb.New(storeEvent.Timestamp()),
		AggregateType:    storeEvent.AggregateType().String(),
		AggregateId:      storeEvent.AggregateID().String(),
		AggregateVersion: &wrapperspb.UInt64Value{Value: storeEvent.AggregateVersion()},
		Data:             storeEvent.Data(),
	}

	return ev, nil
}
