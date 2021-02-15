package util

import (
	"context"

	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	domainRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	eh "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

// SetupQueryHandlerDomain sets up the necessary handlers/projectors/repositories for the query side of es/cqrs.
func SetupQueryHandlerDomain(ctx context.Context, ebConsumer es.EventBusConsumer, esClient esApi.EventStoreClient) (domainRepos.UserRepository, error) {
	// Setup event sourcing
	userRoleBindingRepo := domainRepos.NewUserRoleBindingRepository(esRepos.NewInMemoryProjectionRepository())
	userRepo := domainRepos.NewUserRepository(esRepos.NewInMemoryProjectionRepository(), userRoleBindingRepo)

	err := ebConsumer.AddHandler(ctx,
		es.UseEventHandlerMiddleware(
			eh.NewProjectionRepositoryEventHandler(
				projectors.NewUserProjector(),
				userRepo,
			),
			eh.NewEventStoreReplayMiddleware(esClient).Middleware,
		),
		ebConsumer.Matcher().MatchAggregateType(aggregates.User),
	)
	if err != nil {
		return nil, err
	}

	return userRepo, nil
}

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context) (es.AggregateRepository, error) {
	// ch.NewStoringAggregateHandler()
	return nil, nil
}
