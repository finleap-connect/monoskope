package storage

import (
	"context"
)

// InMemoryEventStore implements an EventStore in memory.
type InMemoryEventStore struct {
	events []Event
}

func NewInMemoryEventStore() Store {
	s := &InMemoryEventStore{
		events: make([]Event, 0),
	}
	return s
}

// Save implements the Save method of the EventStore interface.
func (s *InMemoryEventStore) Save(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return EventStoreError{
			Err: ErrNoEventsToAppend,
		}
	}

	// Validate incoming events and create all event records.
	aggregateID := events[0].AggregateID()
	aggregateType := events[0].AggregateType()
	nextVersion := events[0].AggregateVersion()

	for _, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != aggregateID || event.AggregateType() != aggregateType {
			return EventStoreError{
				Err: ErrInvalidAggregateType,
			}
		}

		// Only accept events that apply to the correct aggregate version.
		if event.AggregateVersion() != nextVersion {
			return EventStoreError{
				Err: ErrIncorrectAggregateVersion,
			}
		}

		for _, se := range s.events {
			if se.AggregateID() == event.AggregateID() && se.AggregateVersion() == event.AggregateVersion() && se.AggregateType() == event.AggregateType() {
				return EventStoreError{
					Err: ErrAggregateVersionAlreadyExists,
				}
			}
		}

		// Increment to checking order of following events.
		nextVersion++
	}

	s.events = append(s.events, events...)

	return nil
}

// Load implements the Load method of the EventStore interface.
func (s *InMemoryEventStore) Load(ctx context.Context, storeQuery *StoreQuery) ([]Event, error) {
	var events []Event

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

// Clear clears the event storage. This is only for testing purposes.
func (s *InMemoryEventStore) clear(ctx context.Context) {
	s.events = make([]Event, 0)
}
