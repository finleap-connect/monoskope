package gateway

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	api_gwauth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type apiServer struct {
	api_gw.UnimplementedGatewayServer
	api_gwauth.UnimplementedAuthServer
	api_common.UnimplementedServiceInformationServiceServer
	// Logger interface
	log logger.Logger
	//
	authConfig  *auth.Config
	authHandler *auth.Handler
}

func NewApiServer(authConfig *auth.Config, authHandler *auth.Handler) *apiServer {
	s := &apiServer{
		log:         logger.WithName("server"),
		authConfig:  authConfig,
		authHandler: authHandler,
	}
	return s
}

func (s *apiServer) GetServiceInformation(e *empty.Empty, stream api_common.ServiceInformationService_GetServiceInformationServer) error {
	err := stream.Send(&api_common.ServiceInformation{
		Name:    "gateway",
		Version: version.Version,
		Commit:  version.Commit,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *apiServer) GetAuthInformation(ctx context.Context, state *api_gwauth.AuthState) (*api_gwauth.AuthInformation, error) {
	url, encodedState, err := s.authHandler.GetAuthCodeURL(state, &auth.AuthCodeURLConfig{
		Scopes:        []string{"offline_access"},
		Clients:       []string{},
		OfflineAccess: true,
	})
	if err != nil {
		return nil, grpcutil.ErrInvalidArgument(err)
	}

	return &api_gwauth.AuthInformation{AuthCodeURL: url, State: encodedState}, nil
}

func (s *apiServer) ExchangeAuthCode(ctx context.Context, code *api_gwauth.AuthCode) (*api_gwauth.AuthResponse, error) {
	token, err := s.authHandler.Exchange(ctx, code.GetCode(), code.CallbackURL)
	if err != nil {
		return nil, err
	}

	claims, err := s.authHandler.VerifyStateAndClaims(ctx, token, code.GetState())
	if err != nil {
		return nil, err
	}

	userInfo := &api_gwauth.AuthResponse{
		AccessToken: &api_gwauth.AccessToken{
			Token:  token.AccessToken,
			Expiry: timestamppb.New(token.Expiry),
		},
		RefreshToken: token.RefreshToken,
		Email:        claims.Email,
	}
	return userInfo, nil
}

func (s *apiServer) RefreshAuth(ctx context.Context, request *api_gwauth.RefreshAuthRequest) (*api_gwauth.AccessToken, error) {
	token, err := s.authHandler.Refresh(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &api_gwauth.AccessToken{
		Token:  token.AccessToken,
		Expiry: timestamppb.New(token.Expiry),
	}, nil
}
