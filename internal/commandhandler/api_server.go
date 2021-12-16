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

package commandhandler

import (
	"context"
	"fmt"
	"time"

	api_domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	api "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	evs "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
	api_domain.UnimplementedCommandHandlerExtensionsServer
	cmdRegistry evs.CommandRegistry
	log         logger.Logger
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(cmdRegistry evs.CommandRegistry) *apiServer {
	return &apiServer{
		cmdRegistry: cmdRegistry,
		log:         logger.WithName("commandhandler-api-server"),
	}
}

func NewServiceClient(ctx context.Context, commandHandlerAddr string) (*grpc.ClientConn, api.CommandHandlerClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(commandHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	return conn, api.NewCommandHandlerClient(conn), nil
}

// Execute implements the API method Execute
func (s *apiServer) Execute(ctx context.Context, command *commands.Command) (*api.CommandReply, error) {
	id, err := uuid.Parse(command.GetId())
	if err != nil {
		return nil, errors.ErrInvalidArgument(fmt.Sprintf("Failed to parse id of command: %s", err.Error()))
	}

	cmd, err := s.cmdRegistry.CreateCommand(id, evs.CommandType(command.Type), command.Data)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	m, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	result, err := s.cmdRegistry.HandleCommand(m.GetContext(), cmd)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	return &api.CommandReply{
		AggregateId: result.Id.String(),
		Version:     result.Version,
	}, nil
}

// GetPermissionModel implements API method GetPermissionModel
func (s *apiServer) GetPermissionModel(ctx context.Context, in *empty.Empty) (*api_domain.PermissionModel, error) {
	permissionModel := &api_domain.PermissionModel{}
	for _, role := range common.Role_name {
		permissionModel.Roles = append(permissionModel.Roles, role)
	}
	for _, scope := range common.Scope_name {
		permissionModel.Scopes = append(permissionModel.Scopes, scope)
	}
	return permissionModel, nil
}

// GetPolicyOverview implements API method GetPolicyOverview
func (s *apiServer) GetPolicyOverview(ctx context.Context, in *empty.Empty) (*api_domain.PolicyOverview, error) {
	policyOverview := &api_domain.PolicyOverview{}
	commandTypes := s.cmdRegistry.GetRegisteredCommandTypes()

	for _, cmdType := range commandTypes {
		command, err := s.cmdRegistry.CreateCommand(uuid.Nil, cmdType, nil)
		if err != nil {
			return nil, err
		}
		policies := command.Policies(ctx)

		for _, p := range policies {
			policyOverview.Policies = append(policyOverview.Policies, &api_domain.Policy{
				Command: cmdType.String(),
				Role:    p.Role().String(),
				Scope:   p.Scope().String(),
			})
		}
	}

	return policyOverview, nil
}
