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

package eventformatter

import (
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"sync"

	"github.com/finleap-connect/monoskope/pkg/logger"
)

type EventFormatterRegistry interface {
	RegisterEventFormatter(es.EventType, EventFormatter) error
	GetEventFormatter(es.EventType) (EventFormatter, error)
}

type eventFormatterRegistry struct {
	log logger.Logger
	mutex sync.RWMutex
	eventFormatters map[es.EventType]EventFormatter
}

// NewEventFormatterRegistry creates a new event formatter registry
func NewEventFormatterRegistry() EventFormatterRegistry {
	return &eventFormatterRegistry{
		log: logger.WithName("event-formatter-registry"),
		eventFormatters: make(map[es.EventType]EventFormatter),
	}
}

// TODO: simplify or feature usecase? (different services register formatters?)

// RegisterEventFormatter registers an event formatter for an event type.
func (r *eventFormatterRegistry) RegisterEventFormatter(eventType es.EventType, eventFormatter EventFormatter) error {
	if eventFormatter == nil {
		r.log.Info("attempt to register invalid event formatter. Event formatter can't be nil")
		return errors.ErrEventFormatterInvalid
	}

	if eventType.String() == "" {
		r.log.Info("attempt to register empty event type")
		return errors.ErrEmptyEventType
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.eventFormatters[eventType]; ok {
		r.log.Info("attempt to register already registered event formatter for event type", "eventType", eventType)
		return errors.ErrEventFormatterForEventTypeAlreadyRegistered
	}
	r.eventFormatters[eventType] = eventFormatter

	r.log.Info("event formatter for event type has been registered.", "eventType", eventType)
	return nil
}

// GetEventFormatter returns the event formatter of the event type registered with RegisterEventFormatter.
func (r *eventFormatterRegistry) GetEventFormatter(eventType es.EventType) (EventFormatter, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if eventFormatter, ok := r.eventFormatters[eventType]; ok {
		return eventFormatter, nil
	}
	r.log.Info("trying to get an event formatter of non-registered event type", "eventType", eventType)
	return nil, errors.ErrEventFormatterForEventTypeNotRegistered
}
