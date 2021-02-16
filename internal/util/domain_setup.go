package util

import (
	"context"

	"github.com/google/uuid"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregateTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	domainHandlers "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
	esm "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esManager "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/manager"
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

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserServiceClient, esClient esApi.EventStoreClient) error {
	// Register aggregates
	aggregateRegistry := es.NewAggregateRegistry()
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserAggregate(id) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserRoleBindingAggregate(id) })

	// Register command handler and middleware
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			esManager.NewAggregateManager(
				aggregateRegistry,
				esClient,
			),
		),
		domainHandlers.NewAuthorizationHandler(
			repos.NewRemoteUserRepository(userService),
		).Middleware,
	)

	// Register commands
	commandRegistry := es.NewCommandRegistry()
	commandRegistry.RegisterCommand(func() es.Command { return commands.NewCreateUserCommand() })
	commandRegistry.RegisterCommand(func() es.Command { return commands.NewCreateUserRoleBindingCommand() })
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)

	return nil
}
