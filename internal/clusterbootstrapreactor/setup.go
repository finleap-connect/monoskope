package clusterbootstrapreactor

import (
	"context"

	eventsourcingApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/certificatemanagement"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/reactors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
)

func SetupClusterBootstrapReactor(ctx context.Context, eventBus eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient, certManager certificatemanagement.CertificateManager) error {
	// Set up JWT signer
	signer := jwt.NewSigner("/etc/reactor/signing.key")

	// Set up middleware
	replayHandler := eventhandler.NewEventStoreReplayEventHandler(esClient)
	reactorEventHandler := eventhandler.NewReactorEventHandler(esClient, reactors.NewClusterBootstrapReactor(signer, certManager))
	//
	reactorHandlerChain := eventsourcing.UseEventHandlerMiddleware(reactorEventHandler, replayHandler.AsMiddleware)

	// Setup matcher for event bus
	clusterCreatedMatcher := eventBus.Matcher().MatchEventType(events.ClusterCreated)

	// Register event handler with event bus
	if err := eventBus.AddHandler(ctx, reactorHandlerChain, clusterCreatedMatcher); err != nil {
		return err
	}

	return nil
}
