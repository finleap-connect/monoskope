package usecases

import (
	"context"
	"fmt"
	"io"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/config"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServerVersionUseCase provides the internal use-case of getting the server version.
type ServerVersionUseCase struct {
	log    logger.Logger
	ctx    context.Context
	config *config.Config
}

func NewServerVersionUseCase(ctx context.Context, config *config.Config) *ServerVersionUseCase {
	useCase := &ServerVersionUseCase{
		log:    logger.WithName("auth-use-case"),
		config: config,
		ctx:    ctx,
	}
	return useCase
}

func (a *ServerVersionUseCase) Run() ([]string, error) {
	conn, err := gateway.CreateGatewayConnecton(a.ctx, a.config.Server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	grpcClient := api_common.NewServiceInformationServiceClient(conn)

	serverInfo, err := grpcClient.GetServiceInformation(a.ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	var serviceInfos []string
	for {
		// Read next
		serverInfo, err := serverInfo.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, err
		}

		// Append
		serviceInfos = append(serviceInfos, fmt.Sprintf(`%s:
		version     : %s
		commit      : %s`,
			serverInfo.GetName(), serverInfo.GetVersion(), serverInfo.GetCommit()))
	}

	return serviceInfos, nil
}
