package queryhandler

import (
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
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
