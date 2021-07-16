package commandhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	api_domain "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
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

	esReply, err := s.cmdRegistry.HandleCommand(m.GetContext(), cmd)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	reply := &api.CommandReply{
		AggregateId: esReply.Id.String(),
		Version:     esReply.Version,
	}
	return reply, nil
}

// GetPermissionModel implements API method GetPermissionModel
func (s *apiServer) GetPermissionModel(ctx context.Context, in *empty.Empty) (*api_domain.PermissionModel, error) {
	permissionModel := &api_domain.PermissionModel{}
	for _, role := range roles.AvailableRoles {
		permissionModel.Roles = append(permissionModel.Roles, role.String())
	}
	for _, scope := range scopes.AvailableScopes {
		permissionModel.Scopes = append(permissionModel.Scopes, scope.String())
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
