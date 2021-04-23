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
	aggregateRegistry.RegisterAggregate(aggregates.NewClusterRegistrationAggregate)
	aggregateRegistry.RegisterAggregate(aggregates.NewClusterAggregate)

	return aggregateManager
}

// setupSuperUsers creates users and rolebindings for all super users
func setupSuperUsers(ctx context.Context, handler es.CommandHandler) error {
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
			cmd, err := es.DefaultCommandRegistry.CreateCommand(userId, commandTypes.CreateUser, data)
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

			cmd, err = es.DefaultCommandRegistry.CreateCommand(uuid.New(), commandTypes.CreateUserRoleBinding, data)
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
func SetupCommandHandlerDomain(ctx context.Context, userService domainApi.UserClient, esClient esApi.EventStoreClient) error {
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

	// Set command handler
	for _, t := range es.DefaultCommandRegistry.GetRegisteredCommandTypes() {
		es.DefaultCommandRegistry.SetHandler(handler, t)
	}

	// Create super users
	cancel := authorizationHandler.BypassAuthorization()
	defer cancel()

	if err := setupSuperUsers(ctx, handler); err != nil {
		return err
	}

	return nil
}
