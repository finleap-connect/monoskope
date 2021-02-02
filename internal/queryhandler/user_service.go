package queryhandler

import (
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/user"
)

// userServiceServer is the implementation of the TenantService API
type userServiceServer struct {
	api.UnimplementedUserServiceServer
	esClient api_es.EventStoreClient

	repo user.ReadOnlyUserRepository
}

// NewUserServiceServer returns a new configured instance of userServiceServer
func NewUserServiceServer(esClient api_es.EventStoreClient, userRepo user.ReadOnlyUserRepository) *userServiceServer {
	return &userServiceServer{
		esClient: esClient,
		repo:     userRepo,
	}
}
