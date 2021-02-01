package queryhandler

import (
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
)

// tenantServiceServer is the implementation of the TenantService API
type tenantServiceServer struct {
	api.UnimplementedTenantServiceServer
	esClient api_es.EventStoreClient
}

// NewTenantServiceServer returns a new configured instance of tenantServiceServer
func NewTenantServiceServer(esClient api_es.EventStoreClient) *tenantServiceServer {
	return &tenantServiceServer{
		esClient: esClient,
	}
}
