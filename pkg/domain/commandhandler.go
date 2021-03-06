// Copyright 2022 Monoskope Authors
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

	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
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
func setupRoleBinding(ctx context.Context, userId uuid.UUID, role, scope string, handler es.CommandHandler) (uuid.UUID, error) {
	data, err := commands.CreateCommandData(&cmdData.CreateUserRoleBindingCommandData{
		UserId: userId.String(),
		Role:   role,
		Scope:  scope,
	})
	if err != nil {
		return uuid.Nil, err
	}

	cmd, err := es.DefaultCommandRegistry.CreateCommand(uuid.New(), commandTypes.CreateUserRoleBinding, data)
	if err != nil {
		return uuid.Nil, err
	}

	reply, err := handler.HandleCommand(ctx, cmd)
	if err != nil && !errors.Is(err, domainErrors.ErrUserRoleBindingAlreadyExists) {
		return uuid.Nil, err
	}

	return reply.Id, nil
}

// setupMockUsers creates mock users for tests
func setupMockUsers(ctx context.Context, handler es.CommandHandler) error {
	createMocks := os.Getenv("CREATE_MOCKS")
	if createMocks != "true" {
		return nil
	}

	for _, mockUser := range mock.TestMockUsers {
		userId, err := setupUser(ctx, mockUser.Name, mockUser.Email, handler)
		if err != nil {
			if errors.Is(err, domainErrors.ErrUserAlreadyExists) {
				return nil
			}
			return err
		}
		mockUser.Id = userId.String()

		for _, mockRole := range mockUser.Roles {
			roleBindingId, err := setupRoleBinding(ctx, userId, mockRole.Role, mockRole.Scope, handler)
			if err != nil {
				return err
			}
			mockRole.Id = roleBindingId.String()
			mockRole.UserId = userId.String()
		}
	}

	return nil
}

// setupSuperUsers creates super users/rolebindings
func setupSuperUsers(ctx context.Context, handler es.CommandHandler) error {
	superUsersEnv := os.Getenv("SUPER_USERS")
	if len(superUsersEnv) < 1 {
		return nil
	}
	superUsers := strings.Split(superUsersEnv, ",")
	if len(superUsers) < 1 {
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

		_, err = setupRoleBinding(ctx, userId, string(roles.Admin), string(scopes.System), handler)
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
	if err := setupMockUsers(ctx, handler); err != nil {
		return err
	}
	return nil
}

// SetupCommandHandlerDomain sets up the necessary handlers/repositories for the command side of es/cqrs.
func SetupCommandHandlerDomain(ctx context.Context, esClient esApi.EventStoreClient) error {
	// Register aggregates
	aggregateManager := registerAggregates(esClient)

	// Create command handler
	handler := es.UseCommandHandlerMiddleware(
		esCommandHandler.NewAggregateHandler(
			aggregateManager,
		),
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

	if err := setupUsers(metadataManager.GetContext(), handler); err != nil {
		return err
	}

	return nil
}
