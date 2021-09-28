// Copyright 2021 Monoskope Authors
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

package usecases

import (
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
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

	return storeQuery, nil
}
