package commandhandler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	commands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
)

// apiServer is the implementation of the CommandHandler API
type apiServer struct {
	api.UnimplementedCommandHandlerServer
	api_common.UnimplementedServiceInformationServiceServer
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer() *apiServer {
	return &apiServer{}
}

// Execute implements the API method Execute
func (s *apiServer) Execute(ctx context.Context, command *commands.Command) (*commands.CommandResult, error) {
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
