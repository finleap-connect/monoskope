package queryhandler

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
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
		return nil, err
	}
	return user.Proto(), nil
}
