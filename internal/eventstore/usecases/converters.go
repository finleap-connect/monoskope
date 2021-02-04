package usecases

import (
	"github.com/google/uuid"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
)

// NewStoreQueryFromProto converts proto api_es.EventFilter to storage.StoreQuery
func NewStoreQueryFromProto(protoFilter *api_es.EventFilter) (*evs.StoreQuery, error) {
	storeQuery := &evs.StoreQuery{}

	if val, ok := protoFilter.GetByAggregate().(*api_es.EventFilter_AggregateId); ok {
		aId, err := uuid.Parse(val.AggregateId.GetValue())
		if err != nil {
			return nil, errors.ErrCouldNotParseAggregateId
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
