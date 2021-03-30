package manager

import (
	"context"

	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
)

// aggregateManager handles storing and loading aggregates in memory.
type aggregateManagerMock struct {
	registry es.AggregateRegistry
	store    es.Store
}

// NewAggregateManager creates a new AggregateHandler which loads/updates Aggregates with the given EventStore.
func NewAggregateManagerMock(aggregateRegistry es.AggregateRegistry) es.AggregateManager {
	return &aggregateManagerMock{
		registry: aggregateRegistry,
		store:    storage.NewInMemoryEventStore(),
	}
}

// Get returns the most recent version of an aggregate.
func (r *aggregateManagerMock) All(ctx context.Context, aggregateType es.AggregateType) ([]es.Aggregate, error) {
	// Retrieve events from store
	eventStream, err := r.store.Load(ctx, &es.StoreQuery{
		AggregateType: &aggregateType,
	})
	if err != nil {
		return nil, err
	}

	// Apply all events gathered from store on aggregate.
	aggregates := make(map[uuid.UUID]es.Aggregate)
	for _, event := range eventStream {
		if event.AggregateType() != aggregateType {
			return nil, errors.ErrInvalidAggregateType
		}

		var aggregate es.Aggregate
		var ok bool
		if aggregate, ok = aggregates[event.AggregateID()]; !ok {
			// Create new empty aggregate of type.
			aggregate, err = r.registry.CreateAggregate(aggregateType, event.AggregateID())
			if err != nil {
				return nil, err
			}
			aggregates[event.AggregateID()] = aggregate
		}

		if err := aggregate.ApplyEvent(event); err != nil {
			return nil, err
		}

		aggregate.IncrementVersion()
	}

	return toAggregateArray(aggregates), nil
}

// Get returns the most recent version of an aggregate.
func (r *aggregateManagerMock) Get(ctx context.Context, aggregateType es.AggregateType, id uuid.UUID) (es.Aggregate, error) {
	// Retrieve events from store
	eventStream, err := r.store.Load(ctx, &es.StoreQuery{
		AggregateId:   &id,
		AggregateType: &aggregateType,
	})
	if err != nil {
		return nil, err
	}

	// Create new empty aggregate of type.
	aggregate, err := r.registry.CreateAggregate(aggregateType, id)
	if err != nil {
		return nil, err
	}

	// Apply all events gathered from store on aggregate.
	for _, event := range eventStream {
		if event.AggregateType() != aggregateType {
			return nil, errors.ErrInvalidAggregateType
		}

		if err := aggregate.ApplyEvent(event); err != nil {
			return nil, err
		}

		aggregate.IncrementVersion()
	}

	return aggregate, nil
}

// Update stores all in-flight events for an aggregate.
func (r *aggregateManagerMock) Update(ctx context.Context, aggregate es.Aggregate) error {
	events := aggregate.Events()

	// Check that there are events in-flight.
	if len(events) == 0 {
		return nil
	}

	// Create stream to send events to store.
	return r.store.Save(ctx, events)
}
