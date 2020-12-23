package usecases

import (
	"errors"

	"github.com/google/uuid"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"
	"google.golang.org/protobuf/encoding/protojson"
)

var ErrCouldNotMarshalEventData = errors.New("could not marshal event data")
var ErrCouldNotParseAggregateId = errors.New("could not parse aggregate id")

// NewEventFromProto converts proto events to storage events
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

func NewStoreQuery(protoFilter *api_es.EventFilter) *storage.StoreQuery {
	// TODO: implement
	panic("not implemented")
}

func NewProtoEvent(storeEvent storage.Event) *api_es.Event {
	// TODO: implement
	panic("not implemented")
}
