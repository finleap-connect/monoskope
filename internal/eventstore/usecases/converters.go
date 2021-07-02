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

	if protoFilter.GetAggregateId() != nil {
		aId, err := uuid.Parse(protoFilter.GetAggregateId().GetValue())
		if err != nil {
			return nil, errors.ErrCouldNotParseAggregateId
		}
		storeQuery.AggregateId = &aId
	}
	if protoFilter.GetAggregateType() != nil {
		aggregateType := es.AggregateType(protoFilter.GetAggregateType().GetValue())
		storeQuery.AggregateType = &aggregateType
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

	storeQuery.ExcludeDeleted = protoFilter.GetExcludeDeleted()

	return storeQuery, nil
}
