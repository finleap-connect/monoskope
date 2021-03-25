package commandhandler

import (
	"context"
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
	"google.golang.org/grpc"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
	api_domain.UnimplementedCommandHandlerExtensionsServer
	cmdRegistry evs.CommandRegistry
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(cmdRegistry evs.CommandRegistry) *apiServer {
	return &apiServer{
		cmdRegistry: cmdRegistry,
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
func (s *apiServer) Execute(ctx context.Context, command *commands.Command) (*empty.Empty, error) {
	id, err := uuid.Parse(command.GetId())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	cmd, err := s.cmdRegistry.CreateCommand(id, evs.CommandType(command.Type), command.Data)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	m, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	err = s.cmdRegistry.HandleCommand(m.GetContext(), cmd)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	return &empty.Empty{}, nil
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
			res := p.Resource()
			sub := p.Subject()

			if res == "" {
				res = "same"
			}

			if sub == "" {
				sub = "self"
			}

			policyOverview.Policies = append(policyOverview.Policies, &api_domain.Policy{
				Command:  cmdType.String(),
				Role:     p.Role().String(),
				Scope:    p.Scope().String(),
				Resource: res,
				Subject:  sub,
			})
		}
	}

	return policyOverview, nil
}
