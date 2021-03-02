package gateway

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	apiCommon "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type apiServer struct {
	api.UnimplementedGatewayServer
	apiCommon.UnimplementedServiceInformationServiceServer
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

func (s *apiServer) GetServiceInformation(e *empty.Empty, stream apiCommon.ServiceInformationService_GetServiceInformationServer) error {
	err := stream.Send(&apiCommon.ServiceInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *apiServer) GetAuthInformation(ctx context.Context, state *api.AuthState) (*api.AuthInformation, error) {
	url, encodedState, err := s.authHandler.GetAuthCodeURL(state, &auth.AuthCodeURLConfig{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid argument: %v", err)
	}

	return &api.AuthInformation{AuthCodeURL: url, State: encodedState}, nil
}

func (s *apiServer) ExchangeAuthCode(ctx context.Context, code *api.AuthCode) (*api.AuthResponse, error) {
	s.log.Info("Authenticating user...")

	token, err := s.authHandler.Exchange(ctx, code.GetCode(), code.CallbackURL)
	if err != nil {
		return nil, err
	}

	claims, err := s.authHandler.VerifyStateAndClaims(ctx, token, code.GetState())
	if err != nil {
		return nil, err
	}

	userInfo := &api.AuthResponse{
		AccessToken:  token.AccessToken,
		IdToken:      auth.GetIDToken(token),
		RefreshToken: token.RefreshToken,
	}

	if !token.Expiry.IsZero() {
		userInfo.Expiry = timestamppb.New(token.Expiry)
	}

	s.log.Info("User authenticated sucessfully.", "User", claims.Email, "Expiry", token.Expiry.String())
	return userInfo, nil
}

func (s *apiServer) RefreshAuth(ctx context.Context, request *api.RefreshAuthRequest) (*api.AuthResponse, error) {
	s.log.Info("Refreshing authentication of user...")

	token, err := s.authHandler.Refresh(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	s.log.Info("Refreshed authentication sucessfully.", "Expiry", token.Expiry.String())
	accessToken := &api.AuthResponse{
		AccessToken:  token.AccessToken,
		IdToken:      auth.GetIDToken(token),
		RefreshToken: token.RefreshToken,
	}
	if !token.Expiry.IsZero() {
		accessToken.Expiry = timestamppb.New(token.Expiry)
	}
	return accessToken, nil
}
