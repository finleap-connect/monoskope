package commandhandler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
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
func (s *apiServer) Execute(ctx context.Context, apiCommand *commands.CommandRequest) (*empty.Empty, error) {
	cmdDetails := apiCommand.GetCommand()

	cmd, err := s.cmdRegistry.CreateCommand(evs.CommandType(cmdDetails.Type), cmdDetails.Data)
	if err != nil {
		return nil, err
	}

	manager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, err
	}

	err = s.cmdRegistry.HandleCommand(manager.SetMetadata(apiCommand.Metadata).GetContext(), cmd)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
