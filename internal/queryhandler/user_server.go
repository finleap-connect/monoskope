package queryhandler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// UserServer is the implementation of the TenantService API
type UserServer struct {
	api.UnimplementedUserServer

	repo repositories.ReadOnlyUserRepository
}

// NewUserServer returns a new configured instance of UserServer
func NewUserServer(userRepo repositories.ReadOnlyUserRepository) *UserServer {
	return &UserServer{
		repo: userRepo,
	}
}

func NewUserClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.UserClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewUserClient(conn), nil
}

// GetById returns the user found by the given id.
func (s *UserServer) GetById(ctx context.Context, userId *wrappers.StringValue) (*projections.User, error) {
	uuid, err := uuid.Parse(userId.Value)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}

	user, err := s.repo.ByUserId(ctx, uuid)
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return user.Proto(), nil
}

// GetByEmail returns the user found by the given email address.
func (s *UserServer) GetByEmail(ctx context.Context, email *wrappers.StringValue) (*projections.User, error) {
	user, err := s.repo.ByEmail(ctx, email.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return user.Proto(), nil
}

func (s *UserServer) GetRoleBindingsById(userId *wrappers.StringValue, stream api.User_GetRoleBindingsByIdServer) error {
	uuid, err := uuid.Parse(userId.Value)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	user, err := s.repo.ByUserId(stream.Context(), uuid)
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

func (s *UserServer) GetAll(request *api.GetAllRequest, stream api.User_GetAllServer) error {
	users, err := s.repo.GetAll(stream.Context(), request.GetIncludeDeleted())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, user := range users {
		err := stream.Send(user.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}
