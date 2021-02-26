package gateway

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	apiCommon "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
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
	authConfig       *auth.Config
	authHandler      *auth.Handler
	cmdHandlerClient apiEs.CommandHandlerClient
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
	url, encodedState, err := s.authHandler.GetAuthCodeURL(state, &auth.AuthCodeURLConfig{
		OfflineAccess: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid argument: %v", err)
	}

	return &api.AuthInformation{AuthCodeURL: url, State: encodedState}, nil
}

func (s *apiServer) ExchangeAuthCode(ctx context.Context, code *api.AuthCode) (*api.AuthResponse, error) {
	token, err := s.authHandler.Exchange(ctx, code.GetCode(), code.CallbackURL)
	if err != nil {
		return nil, err
	}

	claims, err := s.authHandler.VerifyStateAndClaims(ctx, token, code.GetState())
	if err != nil {
		return nil, err
	}

	userInfo := &api.AuthResponse{
		AccessToken: &api.AccessToken{
			Token:  token.AccessToken,
			Expiry: timestamppb.New(token.Expiry),
		},
		RefreshToken: token.RefreshToken,
		Email:        claims.Email,
	}
	return userInfo, nil
}

func (s *apiServer) RefreshAuth(ctx context.Context, request *api.RefreshAuthRequest) (*api.AccessToken, error) {
	token, err := s.authHandler.Refresh(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &api.AccessToken{
		Token:  token.AccessToken,
		Expiry: timestamppb.New(token.Expiry),
	}, nil
}

// Execute implements the API method Execute
func (s *apiServer) Execute(ctx context.Context, command *commands.Command) (*empty.Empty, error) {
	// Get the claims of the authenticated user from the context
	claims, ok := ctx.Value(&auth.Claims{}).(auth.Claims)
	if !ok {
		return nil, grpc.ErrInternal("authentication problem")
	}

	manager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, err
	}

	err = manager.SetUserInformation(&metadata.UserInformation{
		Email:   claims.Email,
		Subject: claims.Subject,
		Issuer:  claims.Issuer,
	})
	if err != nil {
		return nil, err
	}

	// Call command handler to execute
	return s.cmdHandlerClient.Execute(ctx, command)
}
