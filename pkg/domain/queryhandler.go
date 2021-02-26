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
	userEventHandler := es.UseEventHandlerMiddleware(
		esm.NewProjectingEventHandler(
			projectors.NewUserProjector(),
			userRepo,
		),
		esm.NewEventStoreReplayMiddleware(esClient).Middleware,
	)

	// Register event handler with message bus
	err := messageBusConsumer.AddHandler(ctx,
		userEventHandler,
		messageBusConsumer.Matcher().MatchAggregateType(aggregateTypes.User),
	)
	if err != nil {
		return nil, err
	}

	// Setup event handler and middleware
	userRoleBindingEventHandler := es.UseEventHandlerMiddleware(
		esm.NewProjectingEventHandler(
			projectors.NewUserRoleBindingProjector(),
			userRoleBindingRepo,
		),
		esm.NewEventStoreReplayMiddleware(esClient).Middleware,
	)

	// Register event handler with message bus
	err = messageBusConsumer.AddHandler(ctx,
		userRoleBindingEventHandler,
		messageBusConsumer.Matcher().MatchAggregateType(aggregateTypes.UserRoleBinding),
	)
	if err != nil {
		return nil, err
	}

	// Start repo warming
	userRepoWarming := dom.NewRepoWarmingMiddleware(esClient, aggregateTypes.User, userEventHandler)
	if err := userRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}
	userRoleBindingRepoWarming := dom.NewRepoWarmingMiddleware(esClient, aggregateTypes.UserRoleBinding, userRoleBindingEventHandler)
	if err := userRoleBindingRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}

	return userRepo, nil
}
