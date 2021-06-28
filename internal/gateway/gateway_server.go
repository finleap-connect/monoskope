package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gatewayApiServer struct {
	api.UnimplementedGatewayServer
	// Logger interface
	log logger.Logger
	//
	authConfig  *auth.Config
	authHandler *auth.Handler
	userRepo    repositories.ReadOnlyUserRepository
}

func NewGatewayAPIServer(authConfig *auth.Config, authHandler *auth.Handler, userRepo repositories.ReadOnlyUserRepository) api.GatewayServer {
	s := &gatewayApiServer{
		log:         logger.WithName("server"),
		authConfig:  authConfig,
		authHandler: authHandler,
		userRepo:    userRepo,
	}
	return s
}

func (s *gatewayApiServer) GetAuthInformation(ctx context.Context, state *api.AuthState) (*api.AuthInformation, error) {
	url, encodedState, err := s.authHandler.GetAuthCodeURL(state, s.authConfig.Scopes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid argument: %v", err)
	}

	return &api.AuthInformation{AuthCodeUrl: url, State: encodedState}, nil
}

func (s *gatewayApiServer) ExchangeAuthCode(ctx context.Context, code *api.AuthCode) (*api.AuthResponse, error) {
	s.log.Info("Authenticating user...")

	upstreamClaims, err := s.authHandler.Exchange(ctx, code.GetCode(), code.GetState(), code.CallbackUrl)
	if err != nil {
		s.log.Error(err, "User authentication failed.")
		return nil, err
	}

	user, err := s.userRepo.ByEmail(ctx, upstreamClaims.Email)
	if err != nil {
		return nil, err
	}

	signedToken, rawToken, err := s.authHandler.IssueToken(ctx, upstreamClaims, user.Id)
	if err != nil {
		s.log.Error(err, "Issueing token failed.")
		return nil, err
	}

	userInfo := &api.AuthResponse{
		AccessToken: signedToken,
		Expiry:      timestamppb.New(rawToken.Expiry.Time()),
		Username:    user.Name,
	}

	s.log.Info("User authenticated successfully.", "User", upstreamClaims.Email, "Expiry", userInfo.Expiry)

	return userInfo, nil
}
