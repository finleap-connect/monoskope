package domain

import (
	"context"

	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	aggregateTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esm "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

// SetupQueryHandlerDomain sets up the necessary handlers/projectors/repositories for the query side of es/cqrs.
func SetupQueryHandlerDomain(ctx context.Context, esConsumer es.EventBusConsumer, esClient esApi.EventStoreClient) (repos.UserRepository, error) {
	// Setup repositories
	userRepo := repos.NewUserRepository(
		esRepos.NewInMemoryRepository(),
		repos.NewUserRoleBindingRepository(esRepos.NewInMemoryRepository()),
	)

	// Setup event handler and middleware
	err := esConsumer.AddHandler(ctx,
		es.UseEventHandlerMiddleware(
			esm.NewProjectionRepositoryEventHandler(
				projectors.NewUserProjector(),
				userRepo,
			),
			esm.NewEventStoreReplayMiddleware(esClient).Middleware,
		),
		esConsumer.Matcher().MatchAggregateType(aggregateTypes.User),
	)
	if err != nil {
		return nil, err
	}

	return userRepo, nil
}
