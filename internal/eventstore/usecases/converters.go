package usecases

import (
	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

// NewStoreQueryFromProto converts proto esApi.EventFilter to storage.StoreQuery
func NewStoreQueryFromProto(protoFilter *esApi.EventFilter) (*es.StoreQuery, error) {
	storeQuery := &es.StoreQuery{}

	if val, ok := protoFilter.GetByAggregate().(*esApi.EventFilter_AggregateId); ok {
		aId, err := uuid.Parse(val.AggregateId.GetValue())
		if err != nil {
			return nil, errors.ErrCouldNotParseAggregateId
		}
		storeQuery.AggregateId = &aId
	}
	if val, ok := protoFilter.GetByAggregate().(*esApi.EventFilter_AggregateType); ok {
		aType := es.AggregateType(val.AggregateType.GetValue())
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
