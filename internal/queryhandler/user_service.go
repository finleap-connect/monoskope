package queryhandler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// userServiceServer is the implementation of the TenantService API
type userServiceServer struct {
	api.UnimplementedUserServiceServer

	repo repositories.ReadOnlyUserRepository
}

// NewUserServiceServer returns a new configured instance of userServiceServer
func NewUserServiceServer(userRepo repositories.ReadOnlyUserRepository) *userServiceServer {
	return &userServiceServer{
		repo: userRepo,
	}
}

func NewUserServiceClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.UserServiceClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewUserServiceClient(conn), nil
}

// GetById returns the user found by the given id.
func (s *userServiceServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.User, error) {
	user, err := s.repo.ByUserId(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}
	return user.Proto(), nil
}

// GetByEmail returns the user found by the given email address.
func (s *userServiceServer) GetByEmail(ctx context.Context, email *wrappers.StringValue) (*projections.User, error) {
	user, err := s.repo.ByEmail(ctx, email.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return user.Proto(), nil
}

func (s *userServiceServer) GetRoleBindingsById(userId *wrappers.StringValue, stream api.UserService_GetRoleBindingsByIdServer) error {
	user, err := s.repo.ByUserId(stream.Context(), userId.GetValue())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, role := range user.Roles {
		err := stream.Send(role)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}
