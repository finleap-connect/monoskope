// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queryhandler

import (
	"context"
	"time"

	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// tenantServer is the implementation of the TenantService API
type tenantServer struct {
	api.UnimplementedTenantServer

	repoTenant repositories.TenantRepository
	repoUsers  repositories.TenantUserRepository
}

// NewTenantServiceServer returns a new configured instance of tenantServiceServer
func NewTenantServer(tenantRepo repositories.TenantRepository, tenantUserRepo repositories.TenantUserRepository) *tenantServer {
	return &tenantServer{
		repoTenant: tenantRepo,
		repoUsers:  tenantUserRepo,
	}
}

func NewTenantClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.TenantClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithInsecure(queryHandlerAddr).
		WithOpenTelemetry().
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewTenantClient(conn), nil
}

// GetById returns the tenant found by the given id.
func (s *tenantServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.Tenant, error) {
	uuid, err := uuid.Parse(id.GetValue())
	if err != nil {
		return nil, err
	}

	tenant, err := s.repoTenant.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return tenant.Proto(), nil
}

// GetByName returns the tenant found by the given name.
func (s *tenantServer) GetByName(ctx context.Context, name *wrappers.StringValue) (*projections.Tenant, error) {
	tenant, err := s.repoTenant.ByName(ctx, name.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return tenant.Proto(), nil
}

// GetAll returns all tenants.
func (s *tenantServer) GetAll(request *api.GetAllRequest, stream api.Tenant_GetAllServer) error {
	tenants, err := s.repoTenant.AllWith(stream.Context(), request.GetIncludeDeleted())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, t := range tenants {
		err := stream.Send(t.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

// GetUsers returns users belonging to the given tenant id.
func (s *tenantServer) GetUsers(id *wrappers.StringValue, stream api.Tenant_GetUsersServer) error {
	uuid, err := uuid.Parse(id.GetValue())
	if err != nil {
		return err
	}

	tenant, err := s.repoTenant.ById(stream.Context(), uuid)
	if err != nil {
		return err
	}

	// skip deleted
	if tenant.Metadata.GetDeleted() != nil {
		return nil
	}

	users, err := s.repoUsers.GetTenantUsersById(stream.Context(), uuid)
	if err != nil {
		return err
	}

	for _, u := range users {
		err := stream.Send(u.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}
