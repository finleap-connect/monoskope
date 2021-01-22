package commandhandler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	commands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
	api_common.UnimplementedServiceInformationServiceServer
	esClient api_es.EventStoreClient
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(esClient api_es.EventStoreClient) *apiServer {
	return &apiServer{
		esClient: esClient,
	}
}

// Execute implements the API method Execute
func (s *apiServer) Execute(ctx context.Context, command *commands.Command) (*commands.CommandResult, error) {
	cmdDetails := command.GetRequest()

	evsCmd, err := evs.Registry.CreateCommand(evs.CommandType(cmdDetails.Type), cmdDetails.Data)
	if err != nil {
		return nil, err
	}

	err = evs.Registry.HandleCommand(ctx, evsCmd)
	if err != nil {
		return nil, err
	}

	panic("not implemented")
}

// GetServiceInformation implements the API method GetServiceInformation
func (s *apiServer) GetServiceInformation(e *empty.Empty, stream api_common.ServiceInformationService_GetServiceInformationServer) error {
	err := stream.Send(&api_common.ServiceInformation{
		Name:    "commandhandler",
		Version: version.Version,
		Commit:  version.Commit,
	})
	if err != nil {
		return err
	}
	return nil
}
