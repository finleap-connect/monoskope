package eventsourcing

import (
	"context"
	"io"

	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// AggregateStore handles storing and loading Aggregates
type AggregateStore interface {
	// Get returns the most recent version of all aggregate of a given type.
	All(context.Context, AggregateType) ([]Aggregate, error)

	// Get returns the most recent version of an aggregate.
	Get(context.Context, AggregateType, uuid.UUID) (Aggregate, error)

	// Update stores all in-flight events for an aggregate.
	Update(context.Context, Aggregate) error
}

// aggregateStore handles storing and loading aggregates from/to the EventStore.
type aggregateStore struct {
	registry AggregateRegistry
	esClient esApi.EventStoreClient
}

// NewAggregateManager creates a new AggregateHandler which loads/updates Aggregates with the given EventStore.
func NewAggregateManager(aggregateRegistry AggregateRegistry, eventStoreClient esApi.EventStoreClient) AggregateStore {
	return &aggregateStore{
		esClient: eventStoreClient,
		registry: aggregateRegistry,
	}
}

// Get returns the most recent version of an aggregate.
func (r *aggregateStore) All(ctx context.Context, aggregateType AggregateType) ([]Aggregate, error) {
	// Retrieve events from store
	stream, err := r.esClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateType: wrapperspb.String(aggregateType.String()),
	})
	if err != nil {
		return nil, err
	}

	var eventStream []Event
	for {
		// Read next event
		protoEvent, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, err
		}

		event, err := NewEventFromProto(protoEvent)
		if err != nil { // Error converting
			return nil, err
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
	}

	// Apply all events gathered from store on aggregate.
	aggregates := make(map[uuid.UUID]Aggregate)
	for _, event := range eventStream {
		if event.AggregateType() != aggregateType {
			return nil, errors.ErrInvalidAggregateType
		}

		var aggregate Aggregate
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

func toAggregateArray(aggregateMap map[uuid.UUID]Aggregate) []Aggregate {
	var aggregates []Aggregate
	for _, aggregate := range aggregateMap {
		aggregates = append(aggregates, aggregate)
	}
	return aggregates
}

// Get returns the most recent version of an aggregate.
func (r *aggregateStore) Get(ctx context.Context, aggregateType AggregateType, id uuid.UUID) (Aggregate, error) {
	// Retrieve events from store
	stream, err := r.esClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateId:   wrapperspb.String(id.String()),
		AggregateType: wrapperspb.String(aggregateType.String()),
	})
	if err != nil {
		return nil, err
	}

	var eventStream []Event
	for {
		// Read next event
		protoEvent, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, err
		}

		event, err := NewEventFromProto(protoEvent)
		if err != nil { // Error converting
			return nil, err
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
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
func (r *aggregateStore) Update(ctx context.Context, aggregate Aggregate) error {
	events := aggregate.UncommittedEvents()

	// Check that there are events in-flight.
	if len(events) == 0 {
		return nil
	}

	// Create stream to send events to store.
	stream, err := r.esClient.Store(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		// Convert to proto event
		protoEvent := NewProtoFromEvent(event)

		// Send event to store
		err = stream.Send(protoEvent)
		if err != nil {
			return err
		}

		// Apply event on aggregate after successful storage
		err = aggregate.ApplyEvent(event)
		if err != nil {
			return err
		}

		aggregate.IncrementVersion()
	}
	_, err = stream.CloseAndRecv()
	return err
}
