package usecases

import (
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"
)

// NewEventFromProto converts proto events to storage events
func NewEventFromProto(protoEvent *api_es.Event) storage.Event {
	// TODO: implement
	panic("not implemented")
	// return &storage.NewEvent(event.GetSequenceNumber(), storage.EventType(event.GetType())
}

func NewStoreQuery(protoFilter *api_es.EventFilter) *storage.StoreQuery {
	// TODO: implement
	panic("not implemented")
}

func NewProtoEvent(storeEvent storage.Event) *api_es.Event {
	// TODO: implement
	panic("not implemented")
}
