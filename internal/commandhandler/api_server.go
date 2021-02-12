package commandhandler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	commands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
	esClient    api_es.EventStoreClient
	cmdRegistry evs.CommandRegistry
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(esClient api_es.EventStoreClient) *apiServer {
	return &apiServer{
		esClient: esClient,
	}
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

	err = manager.SetUserInformation(&metadata.UserInformation{
		Email:   apiCommand.GetUserMetadata().Email,
		Subject: apiCommand.GetUserMetadata().Subject,
		Issuer:  apiCommand.GetUserMetadata().Issuer,
	})
	if err != nil {
		return nil, err
	}

	err = s.cmdRegistry.HandleCommand(manager.GetContext(), cmd)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
