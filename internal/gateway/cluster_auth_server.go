package gateway

import (
	"context"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/usecases"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type clusterAuthApiServer struct {
	api.UnimplementedClusterAuthServer
	// Logger interface
	log logger.Logger
	//
	signer      jwt.JWTSigner
	userRepo    repositories.ReadOnlyUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
	url         string
	validity    time.Duration
}

func NewClusterAuthAPIServer(
	url string,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
	clusterRepo repositories.ReadOnlyClusterRepository,
	validity time.Duration,
) api.ClusterAuthServer {
	s := &clusterAuthApiServer{
		log:         logger.WithName("server"),
		signer:      signer,
		userRepo:    userRepo,
		clusterRepo: clusterRepo,
		url:         url,
		validity:    validity,
	}
	return s
}

func (s *clusterAuthApiServer) GetAuthToken(ctx context.Context, request *api.ClusterAuthTokenRequest) (*api.ClusterAuthTokenResponse, error) {
	result := new(api.ClusterAuthTokenResponse)
	uc := usecases.NewGetAuthTokenUsecase(request, result, s.signer, s.userRepo, s.clusterRepo, s.url, s.validity)
	err := uc.Run(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
