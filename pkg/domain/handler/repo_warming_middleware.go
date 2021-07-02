package handler

import (
	"context"
	"io"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// NewRepoWarmingMiddleware creates an EventHandler which queries the EventStore to warm up the repository initially.
func WarmUp(ctx context.Context, esClient apiEs.EventStoreClient, aggregateType es.AggregateType, eventHandler es.EventHandler) error {
	log := logger.WithName("repository-warming-middleware").WithValues("aggregateType", aggregateType)
	log.Info("Warming up...")

	// Retrieve events from store
	eventStream, err := esClient.Retrieve(ctx, &apiEs.EventFilter{
		AggregateType: wrapperspb.String(aggregateType.String()),
	})
	if err != nil {
		return err
	}

	appliedEvents := 0
	for {
		// Read next
		protoEvent, err := eventStream.Recv()

		if err != nil {
			if err == io.EOF {
				// End of stream
				break
			} else {
				return err
			}
		}

		// Convert event from api to es
		esEvent, err := es.NewEventFromProto(protoEvent)
		if err != nil {
			return err
		}

		// Let the next handler in the chain handle the event
		err = eventHandler.HandleEvent(ctx, esEvent)
		if err != nil {
			return err
		}

		appliedEvents++
	}

	log.Info("Warmup finished.", "eventsApplied", appliedEvents)

	return nil
}
