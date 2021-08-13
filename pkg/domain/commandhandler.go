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
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/users"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	domainHandlers "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esCommandHandler "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/commandhandler"
)

// registerAggregates registers all aggregates
func registerAggregates(esClient esApi.EventStoreClient) es.AggregateStore {
	aggregateManager := es.NewAggregateManager(es.DefaultAggregateRegistry, esClient)

	// User
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewUserAggregate(aggregateManager) })
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewUserRoleBindingAggregate(aggregateManager) })

	// Tenant
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewTenantAggregate(aggregateManager) })

	// Cluster
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewClusterAggregate(aggregateManager) })

	return aggregateManager
}

// setupUser creates users
func setupUser(ctx context.Context, name, email string, handler es.CommandHandler) (uuid.UUID, error) {
	userId := uuid.New()
	data, err := commands.CreateCommandData(&cmdData.CreateUserCommandData{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return userId, err
	}

	cmd, err := es.DefaultCommandRegistry.CreateCommand(userId, commandTypes.CreateUser, data)
	if err != nil {
		return userId, err
	}

	reply, err := handler.HandleCommand(ctx, cmd)
	if err != nil {
		return uuid.Nil, err
	}
	return reply.Id, nil
}

// setupRoleBinding creates rolebindings
func setupRoleBinding(ctx context.Context, userId uuid.UUID, role, scope string, handler es.CommandHandler) error {
	data, err := commands.CreateCommandData(&cmdData.CreateUserRoleBindingCommandData{
		UserId: userId.String(),
		Role:   role,
		Scope:  scope,
	})
	if err != nil {
		return err
	}

	cmd, err := es.DefaultCommandRegistry.CreateCommand(uuid.New(), commandTypes.CreateUserRoleBinding, data)
	if err != nil {
		return err
	}

	_, err = handler.HandleCommand(ctx, cmd)
	if err != nil && !errors.Is(err, domainErrors.ErrUserRoleBindingAlreadyExists) {
		return err
	}

	return nil
}

// setupSuperUsers creates super users/rolebindings
func setupSuperUsers(ctx context.Context, handler es.CommandHandler) error {
	superUsers := strings.Split(os.Getenv("SUPER_USERS"), ",")
	if len(superUsers) == 0 {
		return nil
	}

	for _, superUser := range superUsers {
		userInfo := strings.Split(superUser, "@")

		userId, err := setupUser(ctx, userInfo[0], superUser, handler)
		if err != nil {
			if errors.Is(err, domainErrors.ErrUserAlreadyExists) {
				return nil
			}
			return err
		}

		err = setupRoleBinding(ctx, userId, roles.Admin.String(), scopes.System.String(), handler)
		if err != nil {
			return err
		}
	}

	return nil
}

// setupUsers creates default users/rolebindings
func setupUsers(ctx context.Context, handler es.CommandHandler) error {
	metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}
	metadataMgr.SetUserInformation(&metadata.UserInformation{
		Id:    users.CommandHandlerUser.ID(),
		Name:  users.CommandHandlerUser.Name,
		Email: users.CommandHandlerUser.Email,
	})
	ctx = metadataMgr.GetContext()

	if err := setupSuperUsers(ctx, handler); err != nil {
		return err
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
	authorizationHandler := domainHandlers.NewUserInformationHandler(userRepo)
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

	// Create default and super users
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	cancel := metadataManager.BypassAuthorization()
	defer cancel()
	if err := setupUsers(metadataManager.GetContext(), handler); err != nil {
		return err
	}

	return nil
}
