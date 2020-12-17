package usecases

import (
	"context"
	"fmt"

	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/monoctl/config"
	"golang.org/x/oauth2"
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

func (a *ServerVersionUseCase) Run() (string, error) {
	conn, err := gateway.CreateGatewayConnecton(a.ctx, a.config.Server, &oauth2.Token{AccessToken: a.config.AuthInformation.Token})
	if err != nil {
		return "", err
	}
	defer conn.Close()
	gwc := api_gw.NewGatewayClient(conn)

	serverInfo, err := gwc.GetServerInfo(a.ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`%s:
		version     : %s
		commit      : %s`,
		a.config.Server, serverInfo.GetVersion(), serverInfo.GetCommit()), nil
}
