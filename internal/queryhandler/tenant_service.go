package queryhandler

import (
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
)

// tenantServiceServer is the implementation of the TenantService API
type tenantServiceServer struct {
	api.UnimplementedTenantServiceServer
}

// NewTenantServiceServer returns a new configured instance of tenantServiceServer
func NewTenantServiceServer() *tenantServiceServer {
	return &tenantServiceServer{}
}
