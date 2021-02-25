package domain

import (
	"context"

	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	aggregateTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	dom "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esm "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

// SetupQueryHandlerDomain sets up the necessary handlers/projectors/repositories for the query side of es/cqrs.
func SetupQueryHandlerDomain(ctx context.Context, messageBusConsumer es.EventBusConsumer, esClient esApi.EventStoreClient) (repos.UserRepository, error) {
	// Setup repositories
	userRoleBindingRepo := repos.NewUserRoleBindingRepository(esRepos.NewInMemoryRepository())
	userRepo := repos.NewUserRepository(esRepos.NewInMemoryRepository(), userRoleBindingRepo)

	// Setup event handler and middleware
	userRepoWarming := dom.NewRepoWarmingMiddleware(esClient, aggregateTypes.User)
	err := messageBusConsumer.AddHandler(ctx,
		es.UseEventHandlerMiddleware(
			esm.NewProjectingEventHandler(
				projectors.NewUserProjector(),
				userRepo,
			),
			userRepoWarming.Middleware,
			esm.NewEventStoreReplayMiddleware(esClient).Middleware,
		),
		messageBusConsumer.Matcher().MatchAggregateType(aggregateTypes.User),
	)
	if err != nil {
		return nil, err
	}

	// Setup event handler and middleware
	userRoleBindingRepoWarming := dom.NewRepoWarmingMiddleware(esClient, aggregateTypes.UserRoleBinding)
	err = messageBusConsumer.AddHandler(ctx,
		es.UseEventHandlerMiddleware(
			esm.NewProjectingEventHandler(
				projectors.NewUserRoleBindingProjector(),
				userRoleBindingRepo,
			),
			userRoleBindingRepoWarming.Middleware,
			esm.NewEventStoreReplayMiddleware(esClient).Middleware,
		),
		messageBusConsumer.Matcher().MatchAggregateType(aggregateTypes.UserRoleBinding),
	)
	if err != nil {
		return nil, err
	}

	// Start repo warming
	if err := userRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}
	if err := userRoleBindingRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}

	return userRepo, nil
}
