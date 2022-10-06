// Copyright 2022 Monoskope Authors
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
	"github.com/finleap-connect/monoskope/internal/telemetry"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/usecase"
	"go.opentelemetry.io/otel/codes"
)

type RetrieveEventsUseCase struct {
	*usecase.UseCaseBase

	store   es.EventStore
	filters []*esApi.EventFilter
	stream  esApi.EventStore_RetrieveServer
	metrics *metrics.EventStoreMetrics
}

func NewRetrieveEventsUseCase(stream esApi.EventStore_RetrieveServer, store es.EventStore, eventFilter *esApi.EventFilter, metrics *metrics.EventStoreMetrics) usecase.UseCase {
	return newRetrieveEventsUseCase(stream, store, &esApi.EventFilters{Filters: []*esApi.EventFilter{eventFilter}}, metrics)
}

func NewRetrieveOrEventsUseCase(stream esApi.EventStore_RetrieveServer, store es.EventStore, eventFilters *esApi.EventFilters, metrics *metrics.EventStoreMetrics) usecase.UseCase {
	return newRetrieveEventsUseCase(stream, store, eventFilters, metrics)
}

// NewRetrieveEventsUseCase creates a new usecase which retrieves all events
// from the store which match the filters
func newRetrieveEventsUseCase(stream esApi.EventStore_RetrieveServer, store es.EventStore, eventFilters *esApi.EventFilters, metrics *metrics.EventStoreMetrics) usecase.UseCase {
	useCase := &RetrieveEventsUseCase{
		UseCaseBase: usecase.NewUseCaseBase("retrieve-events"),
		store:       store,
		filters:     eventFilters.Filters,
		stream:      stream,
		metrics:     metrics,
	}
	return useCase
}

func (u *RetrieveEventsUseCase) Run(ctx context.Context) error {
	ctx, span := telemetry.GetTracer().Start(ctx, "retrieve-events")
	defer span.End()

	// Convert filters
	var sqs []*es.StoreQuery
	for _, filter := range u.filters {
		sq, err := NewStoreQueryFromProto(filter)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		sqs = append(sqs, sq)
	}

	// Retrieve events from Event Store
	u.Log.V(logger.DebugLevel).Info("Retrieving events from the database...")
	span.AddEvent("Retrieving events from the database")
	eventStream, err := u.load(ctx, sqs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	for {
		e, err := eventStream.Receive()
		if err == io.EOF {
			break
		}
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		streamStartTime := time.Now()
		protoEvent := es.NewProtoFromEvent(e)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		err = u.stream.Send(protoEvent)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		// Count retrieved event
		u.metrics.RetrievedTotalCounter.WithLabelValues(protoEvent.Type, protoEvent.AggregateType).Inc()
		u.metrics.RetrievedHistogram.WithLabelValues(protoEvent.Type, protoEvent.AggregateType).Observe(time.Since(streamStartTime).Seconds())
	}

	return nil
}

func (u *RetrieveEventsUseCase) load(ctx context.Context, storeQueries []*es.StoreQuery) (es.EventStreamReceiver, error) {
	if len(storeQueries) == 1 {
		return u.store.Load(ctx, storeQueries[0])
	}
	return u.store.LoadOr(ctx, storeQueries)
}
