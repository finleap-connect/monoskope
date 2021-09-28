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
	"context"
	"io"
	"time"

	"github.com/finleap-connect/monoskope/internal/eventstore/metrics"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/usecase"
)

type RetrieveEventsUseCase struct {
	*usecase.UseCaseBase

	store   es.EventStore
	filter  *esApi.EventFilter
	stream  esApi.EventStore_RetrieveServer
	metrics *metrics.EventStoreMetrics
}

// NewRetrieveEventsUseCase creates a new usecase which retrieves all events
// from the store which match the filter
func NewRetrieveEventsUseCase(stream esApi.EventStore_RetrieveServer, store es.EventStore, filter *esApi.EventFilter, metrics *metrics.EventStoreMetrics) usecase.UseCase {
	useCase := &RetrieveEventsUseCase{
		UseCaseBase: usecase.NewUseCaseBase("retrieve-events"),
		store:       store,
		filter:      filter,
		stream:      stream,
		metrics:     metrics,
	}
	return useCase
}

func (u *RetrieveEventsUseCase) Run(ctx context.Context) error {
	// Convert filter
	sq, err := NewStoreQueryFromProto(u.filter)
	if err != nil {
		return err
	}

	// Retrieve events from Event Store
	u.Log.V(logger.DebugLevel).Info("Retrieving events from the database...")
	eventStream, err := u.store.Load(ctx, sq)
	if err != nil {
		return err
	}

	for {
		e, err := eventStream.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		streamStartTime := time.Now()
		protoEvent := es.NewProtoFromEvent(e)
		if err != nil {
			return err
		}

		err = u.stream.Send(protoEvent)
		if err != nil {
			return err
		}

		// Count retrieved event
		u.metrics.RetrievedTotalCounter.WithLabelValues(protoEvent.Type, protoEvent.AggregateType).Inc()
		u.metrics.RetrievedHistogram.WithLabelValues(protoEvent.Type, protoEvent.AggregateType).Observe(time.Since(streamStartTime).Seconds())
	}

	return nil
}
