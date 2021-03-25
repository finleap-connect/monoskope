package manager

import (
	"context"
	"io"

	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// aggregateManager handles storing and loading aggregates in memory.
type aggregateManager struct {
	registry es.AggregateRegistry
	esClient esApi.EventStoreClient
}

// NewAggregateManager creates a new AggregateHandler which loads/updates Aggregates with the given EventStore.
func NewAggregateManager(aggregateRegistry es.AggregateRegistry, eventStoreClient esApi.EventStoreClient) es.AggregateManager {
	return &aggregateManager{
		esClient: eventStoreClient,
		registry: aggregateRegistry,
	}
}

// Get returns the most recent version of an aggregate.
func (r *aggregateManager) All(ctx context.Context, aggregateType es.AggregateType, excludeDeleted bool) ([]es.Aggregate, error) {
	// Retrieve events from store
	stream, err := r.esClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateType:  wrapperspb.String(aggregateType.String()),
		ExcludeDeleted: excludeDeleted,
	})
	if err != nil {
		return nil, err
	}

	var eventStream []es.Event
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

		event, err := es.NewEventFromProto(protoEvent)
		if err != nil { // Error converting
			return nil, err
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
	}

	// Apply all events gathered from store on aggregate.
	aggregates := make(map[uuid.UUID]es.Aggregate)
	for _, event := range eventStream {
		if event.AggregateType() != aggregateType {
			return nil, errors.ErrInvalidAggregateType
		}

		if aggregate, ok := aggregates[event.AggregateID()]; !ok {
			// Create new empty aggregate of type.
			aggregate, err = r.registry.CreateAggregate(aggregateType, event.AggregateID())
			if err != nil {
				return nil, err
			}
			aggregates[event.AggregateID()] = aggregate
		} else {

			if err := aggregate.ApplyEvent(event); err != nil {
				return nil, err
			}

			aggregate.IncrementVersion()
		}
	}

	return toAggregateArray(aggregates), nil
}

func toAggregateArray(aggregateMap map[uuid.UUID]es.Aggregate) []es.Aggregate {
	var aggregates []es.Aggregate
	for _, aggregate := range aggregateMap {
		aggregates = append(aggregates, aggregate)
	}
	return aggregates
}

// Get returns the most recent version of an aggregate.
func (r *aggregateManager) Get(ctx context.Context, aggregateType es.AggregateType, id uuid.UUID) (es.Aggregate, error) {
	// Retrieve events from store
	stream, err := r.esClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateId:   wrapperspb.String(id.String()),
		AggregateType: wrapperspb.String(aggregateType.String()),
	})
	if err != nil {
		return nil, err
	}

	var eventStream []es.Event
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

		event, err := es.NewEventFromProto(protoEvent)
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
func (r *aggregateManager) Update(ctx context.Context, aggregate es.Aggregate) error {
	events := aggregate.Events()

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
		protoEvent, err := es.NewProtoFromEvent(event)
		if err != nil {
			return err
		}

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
	}
	_, err = stream.CloseAndRecv()
	return err
}
