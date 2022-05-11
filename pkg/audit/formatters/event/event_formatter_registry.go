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

package event

import (
	"sync"

	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"

	"github.com/finleap-connect/monoskope/pkg/logger"
)

// EventFormatterRegistry is the interface definition for an event-formatter registry
type EventFormatterRegistry interface {
	// RegisterEventFormatter registers an event-formatter factory for an event-type.
	RegisterEventFormatter(es.EventType, func(esApi.EventStoreClient) EventFormatter) error
	// CreateEventFormatter returns the event-formatter of the event-type registered with RegisterEventFormatter.
	CreateEventFormatter(esApi.EventStoreClient, es.EventType) (EventFormatter, error)
}

// eventFormatterRegistry is the implementation for the EventFormatterRegistry interface
type eventFormatterRegistry struct {
	log             logger.Logger
	mutex           sync.RWMutex
	eventFormatters map[es.EventType]func(esApi.EventStoreClient) EventFormatter
}

var DefaultEventFormatterRegistry EventFormatterRegistry

func init() {
	DefaultEventFormatterRegistry = NewEventFormatterRegistry()
}

// NewEventFormatterRegistry creates a new event-formatter registry
func NewEventFormatterRegistry() EventFormatterRegistry {
	return &eventFormatterRegistry{
		log:             logger.WithName("event-formatter-registry"),
		eventFormatters: make(map[es.EventType]func(esApi.EventStoreClient) EventFormatter),
	}
}

// RegisterEventFormatter registers an event-formatter factory for an event-type.
// passing an empty event-type will result in errors.ErrEmptyEventType
// passing a nil event-formatter will result in errors.ErrEventFormatterFactoryInvalid
// if an event-formatter factory for the event-type is already registered errors.ErrEventFormatterFactoryForEventTypeAlreadyRegistered is returned
func (r *eventFormatterRegistry) RegisterEventFormatter(eventType es.EventType, factory func(esApi.EventStoreClient) EventFormatter) error {
	if eventType.String() == "" {
		r.log.Info("attempt to register empty event type")
		return errors.ErrEmptyEventType
	}

	if factory == nil {
		r.log.Info("attempt to register invalid event-formatter factory. Event-formatter factory can't be nil")
		return errors.ErrEventFormatterFactoryInvalid
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.eventFormatters[eventType]; ok {
		r.log.Info("attempt to register already registered event-formatter factory for event-type", "eventType", eventType)
		return errors.ErrEventFormatterFactoryForEventTypeAlreadyRegistered
	}
	r.eventFormatters[eventType] = factory

	r.log.V(logger.DebugLevel).Info("event-formatter factory for event-type has been registered.", "eventType", eventType)
	return nil
}

// CreateEventFormatter returns the event-formatter of the event-type registered with RegisterEventFormatter.
// if no event-formatter for the event-type is registered errors.ErrEventFormatterForEventTypeNotRegistered is returned
func (r *eventFormatterRegistry) CreateEventFormatter(esClient esApi.EventStoreClient, eventType es.EventType) (EventFormatter, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if eventFormatter, ok := r.eventFormatters[eventType]; ok {
		return eventFormatter(esClient), nil
	}
	r.log.Info("trying to get an event-formatter of non-registered event-type", "eventType", eventType)
	return nil, errors.ErrEventFormatterForEventTypeNotRegistered
}
