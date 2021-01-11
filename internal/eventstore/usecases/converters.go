package usecases

import (
	"errors"

	"github.com/google/uuid"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ErrCouldNotMarshalEventData is when marshalling event data failed
var ErrCouldNotMarshalEventData = errors.New("could not marshal event data")

// ErrCouldNotUnmarshalEventData is when unmarshalling event data failed
var ErrCouldNotUnmarshalEventData = errors.New("could not unmarshal event data")

// ErrCouldNotParseAggregateId is when an aggregate id could not be parsed as uuid
var ErrCouldNotParseAggregateId = errors.New("could not parse aggregate id")

// NewEventFromProto converts api_es.Event to storage.Event
func NewEventFromProto(protoEvent *api_es.Event) (storage.Event, error) {
	jsonData, err := protojson.Marshal(protoEvent.Data)
	if err != nil {
		return nil, ErrCouldNotMarshalEventData
	}

	aggregateId, err := uuid.Parse(protoEvent.GetAggregateId())
	if err != nil {
		return nil, ErrCouldNotParseAggregateId
	}

	ev := storage.NewEvent(
		storage.EventType(protoEvent.GetType()),
		storage.EventData(jsonData),
		protoEvent.Timestamp.AsTime(),
		storage.AggregateType(protoEvent.GetAggregateType()),
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
		aType := storage.AggregateType(val.AggregateType.GetValue())
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

// NewProtoFromEvent converts storage.Event to api_es.Event
func NewProtoFromEvent(storeEvent storage.Event) (*api_es.Event, error) {
	ev := &api_es.Event{
		Type:             string(storeEvent.EventType()),
		Timestamp:        timestamppb.New(storeEvent.Timestamp()),
		AggregateType:    string(storeEvent.AggregateType()),
		AggregateId:      storeEvent.AggregateID().String(),
		AggregateVersion: &wrapperspb.UInt64Value{Value: storeEvent.AggregateVersion()},
	}

	eventData := &anypb.Any{}
	err := protojson.Unmarshal([]byte(storeEvent.Data()), eventData)
	if err != nil {
		return nil, ErrCouldNotUnmarshalEventData
	}
	ev.Data = eventData

	return ev, nil
}
