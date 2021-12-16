// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"context"
	"errors"
	"os"
	"strings"

	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	domainHandlers "github.com/finleap-connect/monoskope/pkg/domain/handler"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esCommandHandler "github.com/finleap-connect/monoskope/pkg/eventsourcing/commandhandler"
	"github.com/google/uuid"
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

	// TenantClusterBinding
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewTenantClusterBindingAggregate(aggregateManager) })

	// Certificate
	es.DefaultAggregateRegistry.RegisterAggregate(func() es.Aggregate { return aggregates.NewCertificateAggregate(aggregateManager) })

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
func setupRoleBinding(ctx context.Context, userId uuid.UUID, role common.Role, scope common.Scope, handler es.CommandHandler) error {
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

		err = setupRoleBinding(ctx, userId, common.Role_admin, common.Scope_system, handler)
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
