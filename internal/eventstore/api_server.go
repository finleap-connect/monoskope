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

package eventstore

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
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
	s := &apiServer{
		log:     logger.WithName("server"),
		store:   store,
		bus:     bus,
		metrics: metrics.NewEventStoreMetrics(),
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
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter, s.metrics).Run(stream.Context())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}
