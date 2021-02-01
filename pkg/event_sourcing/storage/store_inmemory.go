package storage

import (
	"context"

	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// InMemoryEventStore implements an EventStore in memory.
type InMemoryEventStore struct {
	events []evs.Event
}

func NewInMemoryEventStore() evs.Store {
	s := &InMemoryEventStore{
		events: make([]evs.Event, 0),
	}
	return s
}

// Connect does nothing
func (b *InMemoryEventStore) Connect(ctx context.Context) error {
	return nil
}

// Save implements the Save method of the EventStore interface.
func (s *InMemoryEventStore) Save(ctx context.Context, events []evs.Event) error {
	if len(events) == 0 {
		return evs.NewEventStoreError(evs.ErrNoEventsToAppend, nil)
	}

	// Validate incoming events and create all event records.
	aggregateID := events[0].AggregateID()
	aggregateType := events[0].AggregateType()
	nextVersion := events[0].AggregateVersion()

	for _, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != aggregateID || event.AggregateType() != aggregateType {
			return evs.NewEventStoreError(evs.ErrInvalidAggregateType, nil)
		}

		// Only accept events that apply to the correct aggregate version.
		if event.AggregateVersion() != nextVersion {
			return evs.NewEventStoreError(evs.ErrIncorrectAggregateVersion, nil)
		}

		for _, se := range s.events {
			if se.AggregateID() == event.AggregateID() && se.AggregateVersion() == event.AggregateVersion() && se.AggregateType() == event.AggregateType() {
				return evs.NewEventStoreError(evs.ErrAggregateVersionAlreadyExists, nil)
			}
		}

		// Increment to checking order of following events.
		nextVersion++
	}

	s.events = append(s.events, events...)

	return nil
}

// Load implements the Load method of the EventStore interface.
func (s *InMemoryEventStore) Load(ctx context.Context, storeQuery *evs.StoreQuery) ([]evs.Event, error) {
	var events []evs.Event

	for _, ev := range s.events {
		if storeQuery.AggregateId != nil && ev.AggregateID() != *storeQuery.AggregateId {
			continue
		} else if storeQuery.AggregateType != nil && ev.AggregateType() != *storeQuery.AggregateType {
			continue
		}

		if storeQuery.MinVersion != nil && ev.AggregateVersion() < *storeQuery.MinVersion {
			continue
		}
		if storeQuery.MaxVersion != nil && ev.AggregateVersion() > *storeQuery.MaxVersion {
			continue
		}

		if storeQuery.MinTimestamp != nil &&
			(ev.Timestamp() != *storeQuery.MinTimestamp || storeQuery.MinTimestamp.Before(ev.Timestamp())) {
			continue
		}
		if storeQuery.MaxTimestamp != nil &&
			(ev.Timestamp() != *storeQuery.MaxTimestamp || storeQuery.MaxTimestamp.After(ev.Timestamp())) {
			continue
		}

		events = append(events, ev)
	}

	return events, nil
}

func (s *InMemoryEventStore) Close() error {
	return nil
}

// Clear clears the event storage. This is only for testing purposes.
func (s *InMemoryEventStore) clear(ctx context.Context) {
	s.events = make([]evs.Event, 0)
}
