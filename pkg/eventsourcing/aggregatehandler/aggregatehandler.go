package aggregatehandler

import (
	"context"
	"io"

	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// aggregateHandler handles storing and loading aggregates in memory.
type aggregateHandler struct {
	eventStoreClient esApi.EventStoreClient
}

// NewAggregateHandler creates a new AggregateHandler which loads/updates Aggregates with the given EventStore.
func NewAggregateHandler(eventStoreClient esApi.EventStoreClient) es.AggregateHandler {
	return &aggregateHandler{
		eventStoreClient: eventStoreClient,
	}
}

// Get returns the most recent version of an aggregate.
func (r *aggregateHandler) Get(ctx context.Context, aggregateType es.AggregateType, id uuid.UUID) (es.Aggregate, error) {
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
func (r *aggregateHandler) Update(context.Context, es.Aggregate) error {
	panic("not implemented")
}
