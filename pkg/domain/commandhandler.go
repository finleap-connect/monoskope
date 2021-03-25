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

func RegisterCommands(superusers ...string) es.CommandRegistry {
	commandRegistry := es.NewCommandRegistry()
	commandRegistry.RegisterCommand(func(id uuid.UUID) es.Command { return commands.NewCreateUserCommand(id) })
	commandRegistry.RegisterCommand(func(id uuid.UUID) es.Command {
		cmd := commands.NewCreateUserRoleBindingCommand(id)
		cmd.DeclareSuperusers(superusers)
		return cmd
	})
	commandRegistry.RegisterCommand(func(id uuid.UUID) es.Command { return commands.NewCreateTenantCommand(id) })
	commandRegistry.RegisterCommand(func(id uuid.UUID) es.Command { return commands.NewUpdateTenantCommand(id) })
	commandRegistry.RegisterCommand(func(id uuid.UUID) es.Command { return commands.NewDeleteTenantCommand(id) })

	return commandRegistry
}

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserServiceClient, tenantService domainApi.TenantServiceClient, esClient esApi.EventStoreClient, superusers ...string) (es.CommandRegistry, error) {
	// Setup repositories
	userRepo := repos.NewRemoteUserRepository(userService)
	tenantRepo := repos.NewRemoteTenantRepository(tenantService)

	// Register aggregates
	aggregateRegistry := es.NewAggregateRegistry()
	aggregateManager := esManager.NewAggregateManager(
		aggregateRegistry,
		esClient,
	)

	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserAggregate(id, aggregateManager) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserRoleBindingAggregate(id, aggregateManager) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewTenantAggregate(id, tenantRepo) })

	// Register command handler and middleware
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			aggregateManager,
		),
		domainHandlers.NewAuthorizationHandler(userRepo).Middleware,
	)

	// Register commands
	commandRegistry := RegisterCommands(superusers...)
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)
	commandRegistry.SetHandler(handler, commandTypes.CreateTenant)
	commandRegistry.SetHandler(handler, commandTypes.UpdateTenant)
	commandRegistry.SetHandler(handler, commandTypes.DeleteTenant)

	return commandRegistry, nil
}
