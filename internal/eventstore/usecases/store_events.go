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

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/internal/eventstore/metrics"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	MAX_BACKOFF_PUBLISH = 1 * time.Minute
)

type StoreEventsUseCase struct {
	*usecase.UseCaseBase

	store   es.EventStore
	bus     es.EventBusPublisher
	stream  esApi.EventStore_StoreServer
	metrics *metrics.EventStoreMetrics
}

// NewStoreEventsUseCase creates a new usecase which stores all events in the store
// and broadcasts these events via the message bus
func NewStoreEventsUseCase(stream esApi.EventStore_StoreServer, store es.EventStore, bus es.EventBusPublisher, metrics *metrics.EventStoreMetrics) usecase.UseCase {
	useCase := &StoreEventsUseCase{
		UseCaseBase: usecase.NewUseCaseBase("store-events"),
		store:       store,
		bus:         bus,
		stream:      stream,
		metrics:     metrics,
	}
	return useCase
}

func (u *StoreEventsUseCase) Run(ctx context.Context) error {
	for {
		startTime := time.Now()

		// Read next event
		event, err := u.stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}

		if err != nil { // Some other error
			return errors.TranslateToGrpcError(err)
		}

		// Count transmitted event
		u.metrics.TransmittedTotalCounter.WithLabelValues(event.Type, event.AggregateType).Inc()

		// Convert from proto events to storage events
		ev, err := es.NewEventFromProto(event)
		if err != nil {
			return err
		}

		// Store events in database
		u.Log.V(logger.DebugLevel).Info("Saving events in the store...")
		if err := u.store.Save(ctx, []es.Event{ev}); err != nil {
			return err
		}

		// Count successfully stored event
		u.metrics.StoredTotalCounter.WithLabelValues(event.Type, event.AggregateType).Inc()

		// Send events to message bus
		u.Log.V(logger.DebugLevel).Info("Sending events to the message bus...")

		params := backoff.NewExponentialBackOff()
		params.MaxElapsedTime = MAX_BACKOFF_PUBLISH

		err = backoff.Retry(func() error {
			return u.bus.PublishEvent(ctx, ev)
		}, params)
		if err != nil {
			return err
		}
		u.metrics.StoredHistogram.WithLabelValues(event.Type, event.AggregateType).Observe(time.Since(startTime).Seconds())
	}

	return u.stream.SendAndClose(&emptypb.Empty{})
}
