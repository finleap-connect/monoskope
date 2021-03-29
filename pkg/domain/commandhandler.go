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
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
	esManager "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/manager"
)

func RegisterCommands() es.CommandRegistry {
	commandRegistry := es.NewCommandRegistry()

	// User
	commandRegistry.RegisterCommand(commands.NewCreateUserCommand)
	commandRegistry.RegisterCommand(commands.NewCreateUserRoleBindingCommand)

	// Tenant
	commandRegistry.RegisterCommand(commands.NewCreateTenantCommand)
	commandRegistry.RegisterCommand(commands.NewUpdateTenantCommand)
	commandRegistry.RegisterCommand(commands.NewDeleteTenantCommand)

	return commandRegistry
}

func SetupSuperUsers(ctx context.Context, commandRegistry es.CommandRegistry, handler es.CommandHandler) error {
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
				Issuer: "system",
			})
			ctx := metadataMgr.GetContext()

			userId := uuid.New()
			data, err := commands.CreateCommandData(&cmdData.CreateUserCommandData{
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
			if err != nil && !errors.Is(err, domainErrors.ErrUserAlreadyExists) {
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
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserServiceClient, esClient esApi.EventStoreClient) (es.CommandRegistry, error) {
	// Setup repositories
	userRepo := repos.NewRemoteUserRepository(userService)

	// Register aggregates
	aggregateRegistry := es.NewAggregateRegistry()
	aggregateManager := esManager.NewAggregateManager(
		aggregateRegistry,
		esClient,
	)

	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserAggregate(id, aggregateManager) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewUserRoleBindingAggregate(id, aggregateManager) })
	aggregateRegistry.RegisterAggregate(func(id uuid.UUID) es.Aggregate { return aggregates.NewTenantAggregate(id, aggregateManager) })

	// Register command handler and middleware
	authorizationHandler := domainHandlers.NewAuthorizationHandler(userRepo)
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			aggregateManager,
		),
		authorizationHandler.Middleware,
	)

	// Register commands
	commandRegistry := RegisterCommands()
	commandRegistry.SetHandler(handler, commandTypes.CreateUser)
	commandRegistry.SetHandler(handler, commandTypes.CreateUserRoleBinding)
	commandRegistry.SetHandler(handler, commandTypes.DeleteUserRoleBinding)
	commandRegistry.SetHandler(handler, commandTypes.CreateTenant)
	commandRegistry.SetHandler(handler, commandTypes.UpdateTenant)
	commandRegistry.SetHandler(handler, commandTypes.DeleteTenant)

	// Create super users
	cancelBypass := authorizationHandler.BypassAuthorization()
	defer cancelBypass()
	if err := SetupSuperUsers(ctx, commandRegistry, handler); err != nil {
		return nil, err
	}

	return commandRegistry, nil
}
