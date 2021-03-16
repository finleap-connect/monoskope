package domain

import (
	"context"

	"github.com/google/uuid"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	domainHandlers "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
	esManager "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/manager"
)

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserServiceClient, tenantService domainApi.TenantServiceClient, esClient esApi.EventStoreClient, superusers ...string) (es.CommandRegistry, error) {
	// Setup repositories
	userRepo := repos.NewRemoteUserRepository(userService)
	tenantRepo := repos.NewRemoteTenantRepository(tenantService)

	// Register aggregates
	aggregateRegistry := es.NewAggregateRegistry()
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserAggregate(id, userRepo) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserRoleBindingAggregate(id) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewTenantAggregate(id, tenantRepo) })

	// Register command handler and middleware
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			esManager.NewAggregateManager(
				aggregateRegistry,
				esClient,
			),
		),
		domainHandlers.NewAuthorizationHandler(userRepo).Middleware,
	)

	// Register commands
	commandRegistry := es.NewCommandRegistry()
	commandRegistry.RegisterCommand(func() es.Command { return commands.NewCreateUserCommand() })
	commandRegistry.RegisterCommand(func() es.Command {
		cmd := commands.NewCreateUserRoleBindingCommand()
		cmd.DeclareSuperusers(superusers)
		return cmd
	})
	commandRegistry.RegisterCommand(func() es.Command { return commands.NewCreateTenantCommand() })
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)
	commandRegistry.SetHandler(handler, commandTypes.CreateTenant)

	return commandRegistry, nil
}
