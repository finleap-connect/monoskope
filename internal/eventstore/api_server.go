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

package eventstore

import (
	"fmt"

	"github.com/finleap-connect/monoskope/internal/eventstore/metrics"
	"github.com/finleap-connect/monoskope/internal/eventstore/usecases"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
)

// apiServer is the implementation of the EventStore API
type apiServer struct {
	esApi.UnimplementedEventStoreServer
	// Logger interface
	log     logger.Logger
	store   es.EventStore
	bus     es.EventBusPublisher
	metrics *metrics.EventStoreMetrics
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(store es.EventStore, bus es.EventBusPublisher) *apiServer {
	m, err := metrics.NewEventStoreMetrics()
	if err != nil {
		panic(fmt.Errorf("Error setting up metrics server: %w", err))
	}

	s := &apiServer{
		log:     logger.WithName("server"),
		store:   store,
		bus:     bus,
		metrics: m,
	}

	return s
}

// Store implements the API method for storing events
func (s *apiServer) Store(stream esApi.EventStore_StoreServer) error {
	// Perform the use case for storing events
	if err := usecases.NewStoreEventsUseCase(stream, s.store, s.bus, s.metrics).Run(stream.Context()); err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}

// Retrieve implements the API method for retrieving events from the store
func (s *apiServer) Retrieve(filter *esApi.EventFilter, stream esApi.EventStore_RetrieveServer) error {
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter, s.metrics).Run(stream.Context())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}

// RetrieveOr implements the API method for retrieving events with the logical OR from the store
func (s *apiServer) RetrieveOr(filters *esApi.EventFilters, stream esApi.EventStore_RetrieveOrServer) error {
	err := usecases.NewRetrieveOrEventsUseCase(stream, s.store, filters, s.metrics).Run(stream.Context())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}
