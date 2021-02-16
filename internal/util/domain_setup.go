package util

import (
	"context"

	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	domainHandlers "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	domainRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
	esMiddleware "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esManager "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/manager"
	esRepos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

// SetupQueryHandlerDomain sets up the necessary handlers/projectors/repositories for the query side of es/cqrs.
func SetupQueryHandlerDomain(ctx context.Context, esConsumer es.EventBusConsumer, esClient esApi.EventStoreClient) (domainRepos.UserRepository, error) {
	// Setup event sourcing
	userRoleBindingRepo := domainRepos.NewUserRoleBindingRepository(esRepos.NewInMemoryRepository())
	userRepo := domainRepos.NewUserRepository(esRepos.NewInMemoryRepository(), userRoleBindingRepo)

	err := esConsumer.AddHandler(ctx,
		es.UseEventHandlerMiddleware(
			esMiddleware.NewProjectionRepositoryEventHandler(
				projectors.NewUserProjector(),
				userRepo,
			),
			esMiddleware.NewEventStoreReplayMiddleware(esClient).Middleware,
		),
		esConsumer.Matcher().MatchAggregateType(aggregates.User),
	)
	if err != nil {
		return nil, err
	}

	return userRepo, nil
}

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserServiceClient, esClient esApi.EventStoreClient) error {
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			esManager.NewAggregateManager(
				es.NewAggregateRegistry(),
				esClient,
			),
		),
		domainHandlers.NewAuthorizationHandler(
			domainRepos.NewRemoteUserRepository(userService),
		).Middleware,
	)

	commandRegistry := es.NewCommandRegistry()
	commandRegistry.RegisterCommand(func() es.Command { return &commands.CreateUserCommand{} })
	commandRegistry.RegisterCommand(func() es.Command { return &commands.CreateUserRoleBindingCommand{} })
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)

	return nil
}
