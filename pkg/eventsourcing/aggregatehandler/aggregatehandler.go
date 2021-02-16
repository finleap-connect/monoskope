package aggregatehandler

import (
	"context"
	"io"

	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// aggregateManager handles storing and loading aggregates in memory.
type aggregateManager struct {
	eventStoreClient esApi.EventStoreClient
}

// NewAggregateManager creates a new AggregateHandler which loads/updates Aggregates with the given EventStore.
func NewAggregateManager(eventStoreClient esApi.EventStoreClient) es.AggregateManager {
	return &aggregateManager{
		eventStoreClient: eventStoreClient,
	}
}

// Get returns the most recent version of an aggregate.
func (r *aggregateManager) Get(ctx context.Context, aggregateType es.AggregateType, id uuid.UUID) (es.Aggregate, error) {
	stream, err := r.eventStoreClient.Retrieve(ctx, &esApi.EventFilter{
		AggregateId:   wrapperspb.String(id.String()),
		AggregateType: wrapperspb.String(aggregateType.String()),
	})
	if err != nil {
		return nil, err
	}

	var eventStream []*esApi.Event
	for {
		// Read next event
		event, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, err
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
	}

	panic("not implemented")
}

// Update stores all in-flight events for an aggregate.
func (r *aggregateManager) Update(context.Context, es.Aggregate) error {
	panic("not implemented")
}
