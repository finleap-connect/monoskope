package handler

import (
	"context"
	"io"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type repoWarmingReplayMiddleware struct {
	log                logger.Logger
	esClient           apiEs.EventStoreClient
	aggregateType      es.AggregateType
	nextHandlerInChain es.EventHandler
}

// NewRepoWarmingMiddleware creates an EventHandler which queryies the EventStore to warm up the repository initially.
func NewRepoWarmingMiddleware(esClient apiEs.EventStoreClient, aggregateType es.AggregateType) *repoWarmingReplayMiddleware {
	return &repoWarmingReplayMiddleware{
		log:           logger.WithName("repository-warming-middleware").WithValues("aggregateType", aggregateType),
		esClient:      esClient,
		aggregateType: aggregateType,
	}
}

func (m *repoWarmingReplayMiddleware) Middleware(h es.EventHandler) es.EventHandler {
	m.nextHandlerInChain = h
	return m
}

func (m *repoWarmingReplayMiddleware) WarmUp(ctx context.Context) error {
	m.log.Info("Warming up...")

	// Retrieve events from store
	eventStream, err := m.esClient.Retrieve(ctx, &apiEs.EventFilter{
		AggregateType: wrapperspb.String(m.aggregateType.String()),
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
		err = m.nextHandlerInChain.HandleEvent(ctx, esEvent)
		if err != nil {
			return err
		}

		appliedEvents++
	}

	m.log.Info("Warmup finished.", "eventsApplied", appliedEvents)

	return nil
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *repoWarmingReplayMiddleware) HandleEvent(ctx context.Context, event es.Event) error {
	return m.nextHandlerInChain.HandleEvent(ctx, event)
}
