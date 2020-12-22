package storage

import (
	"errors"
	"fmt"
	"sync"
)

// ErrEventDataNotRegistered is when no event data factory was registered.
var ErrEventDataNotRegistered = errors.New("event data not registered")

// RegisterEventData registers an event data factory for a type. The factory is
// used to create concrete event data structs when loading from the database.
//
// An example would be:
//     RegisterEventData(MyEventType, func() Event { return &MyEventData{} })
func RegisterEventData(eventType EventType, factory func() EventData) {
	if eventType == EventType("") {
		panic("attempt to register empty event type")
	}

	eventDataFactoriesMu.Lock()
	defer eventDataFactoriesMu.Unlock()
	if _, ok := eventDataFactories[eventType]; ok {
		panic(fmt.Sprintf("registering duplicate types for %q", eventType))
	}
	eventDataFactories[eventType] = factory
}

// UnregisterEventData removes the registration of the event data factory for
// a type. This is mainly useful in mainenance situations where the event data
// needs to be switched in a migrations.
func UnregisterEventData(eventType EventType) {
	if eventType == EventType("") {
		panic("attempt to unregister empty event type")
	}

	eventDataFactoriesMu.Lock()
	defer eventDataFactoriesMu.Unlock()
	if _, ok := eventDataFactories[eventType]; !ok {
		panic(fmt.Sprintf("unregister of non-registered type %q", eventType))
	}
	delete(eventDataFactories, eventType)
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
