package storage

import (
	"errors"
	"sync"
)

// ErrEventDataNotRegistered is when no event data factory was registered.
var ErrEventDataNotRegistered = errors.New("event data not registered")

// ErrEventTypeEmpty is when one tries to register an empty event type.
var ErrEventTypeEmpty = errors.New("attempt to register empty event type")

// ErrEventTypeDuplicate is when one tries to register an empty event type.
var ErrEventTypeDuplicate = errors.New("attempt to register event type which is already registered")

// ErrEventTypeUnknown is when one tries to unregister an event type which hasn't been registered before.
var ErrEventTypeUnknown = errors.New("unregister of non-registered type")

const EmptyEventType = EventType("")

// RegisterEventData registers an event data factory for a type. The factory is
// used to create concrete event data structs when loading from the database.
//
// An example would be:
//     RegisterEventData(MyEventType, func() Event { return &MyEventData{} })
func RegisterEventData(eventType EventType, factory func() EventData) error {
	if eventType == EmptyEventType {
		return ErrEventTypeEmpty
	}

	eventDataFactoriesMu.Lock()
	defer eventDataFactoriesMu.Unlock()
	if _, ok := eventDataFactories[eventType]; ok {
		return ErrEventTypeDuplicate
	}
	eventDataFactories[eventType] = factory

	return nil
}

// UnregisterEventData removes the registration of the event data factory for
// a type. This is mainly useful in mainenance situations where the event data
// needs to be switched in a migrations.
func UnregisterEventData(eventType EventType) error {
	if eventType == EmptyEventType {
		return ErrEventTypeEmpty
	}

	eventDataFactoriesMu.Lock()
	defer eventDataFactoriesMu.Unlock()
	if _, ok := eventDataFactories[eventType]; !ok {
		return ErrEventTypeUnknown
	}
	delete(eventDataFactories, eventType)

	return nil
}

// CreateEventData creates an event data of a type using the factory registered
// with RegisterEventData.
func CreateEventData(eventType EventType) (EventData, error) {
	eventDataFactoriesMu.RLock()
	defer eventDataFactoriesMu.RUnlock()
	if factory, ok := eventDataFactories[eventType]; ok {
		return factory(), nil
	}
	return nil, ErrEventDataNotRegistered
}

var eventDataFactories = make(map[EventType]func() EventData)
var eventDataFactoriesMu sync.RWMutex
