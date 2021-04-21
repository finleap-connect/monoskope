package domain

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	domainHandlers "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
	esManager "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/manager"
)

// RegisterCommands registers all commands available
func RegisterCommands() es.CommandRegistry {
	commandRegistry := es.NewCommandRegistry()

	// User
	commandRegistry.RegisterCommand(commands.NewCreateUserCommand)
	commandRegistry.RegisterCommand(commands.NewCreateUserRoleBindingCommand)
	commandRegistry.RegisterCommand(commands.NewDeleteUserRoleBindingCommand)

	// Tenant
	commandRegistry.RegisterCommand(commands.NewCreateTenantCommand)
	commandRegistry.RegisterCommand(commands.NewUpdateTenantCommand)
	commandRegistry.RegisterCommand(commands.NewDeleteTenantCommand)

	// Cluster
	commandRegistry.RegisterCommand(commands.NewRequestClusterRegistrationCommand)
	commandRegistry.RegisterCommand(commands.NewDeleteClusterCommand)

	return commandRegistry
}

// registerCommandsWithHandler registers all commands available and sets the given commandhandler
func registerCommandsWithHandler(handler es.CommandHandler) es.CommandRegistry {
	commandRegistry := RegisterCommands()

	// User
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)
	commandRegistry.SetHandler(handler, commandTypes.DeleteUserRoleBinding)

	// Tenant
	commandRegistry.SetHandler(handler, commandTypes.CreateTenant)
	commandRegistry.SetHandler(handler, commandTypes.UpdateTenant)
	commandRegistry.SetHandler(handler, commandTypes.DeleteTenant)

	// Cluster
	commandRegistry.SetHandler(handler, commandTypes.RequestClusterRegistration)
	commandRegistry.SetHandler(handler, commandTypes.DeleteCluster)

	return commandRegistry
}

// registerAggregates registers all aggregates
func registerAggregates(esClient esApi.EventStoreClient) es.AggregateManager {
	aggregateRegistry := es.NewAggregateRegistry()
	aggregateManager := esManager.NewAggregateManager(
		aggregateRegistry,
		esClient,
	)

	// User
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserAggregate(id, aggregateManager) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserRoleBindingAggregate(id, aggregateManager) })

	// Tenant
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewTenantAggregate(id, aggregateManager) })

	// Cluster
	aggregateRegistry.RegisterAggregate(aggregates.NewClusterAggregate)

	return aggregateManager
}

// setupSuperUsers creates users and rolebindings for all super users
func setupSuperUsers(ctx context.Context, commandRegistry es.CommandRegistry, handler es.CommandHandler) error {
	if superUsers := strings.Split(os.Getenv("SUPERUSERS"), ","); len(superUsers) != 0 {
		for _, superUser := range superUsers {
			userInfo := strings.Split(superUser, "@")
			metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
			if err != nil {
				return err
			}
			metadataMgr.SetUserInformation(&metadata.UserInformation{
				Name:   userInfo[0],
				Email:  superUser,
				Issuer: "commandhandler",
			})
			ctx := metadataMgr.GetContext()

			userId := uuid.New()
			data, err := commands.CreateCommandData(&cmdData.CreateUserCommandData{
				Name:  userInfo[0],
				Email: superUser,
			})
			if err != nil {
				return err
			}
			cmd, err := commandRegistry.CreateCommand(userId, commandTypes.CreateUser, data)
			if err != nil {
				return err
			}

			err = handler.HandleCommand(ctx, cmd)
			if err != nil {
				if errors.Is(err, domainErrors.ErrUserAlreadyExists) {
					continue
				}
				return err
			}

			data, err = commands.CreateCommandData(&cmdData.CreateUserRoleBindingCommandData{
				UserId: userId.String(),
				Role:   roles.Admin.String(),
				Scope:  scopes.System.String(),
			})
			if err != nil {
				return err
			}

			cmd, err = commandRegistry.CreateCommand(uuid.New(), commandTypes.CreateUserRoleBinding, data)
			if err != nil {
				return err
			}

			err = handler.HandleCommand(ctx, cmd)
			if err != nil && !errors.Is(err, domainErrors.ErrUserRoleBindingAlreadyExists) {
				return err
			}
		}
	}
	return nil
}

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserClient, esClient esApi.EventStoreClient) (es.CommandRegistry, error) {
	// Register aggregates
	aggregateManager := registerAggregates(esClient)

	// Setup repositories
	userRepo := repositories.NewRemoteUserRepository(userService)

	// Create command handler
	authorizationHandler := domainHandlers.NewAuthorizationHandler(userRepo)
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			aggregateManager,
		),
		authorizationHandler.Middleware,
	)

	// Register commands
	commandRegistry := registerCommandsWithHandler(handler)

	// Create super users
	cancel := authorizationHandler.BypassAuthorization()
	defer cancel()

	if err := setupSuperUsers(ctx, commandRegistry, handler); err != nil {
		return nil, err
	}

	return commandRegistry, nil
}
