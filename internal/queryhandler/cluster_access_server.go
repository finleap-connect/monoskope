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
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// clusterAccessServer is the implementation of the ClusterAccessService API
type clusterAccessServer struct {
	api.UnimplementedClusterAccessServer
	clusterAccessRepo        repositories.ClusterAccessRepository
	tenantClusterBindingRepo repositories.TenantClusterBindingRepository
}

// NewClusterServiceServer returns a new configured instance of clusterServiceServer
func NewClusterAccessServer(clusterAccessRepo repositories.ClusterAccessRepository, tenantClusterBindingRepo repositories.TenantClusterBindingRepository) *clusterAccessServer {
	return &clusterAccessServer{
		clusterAccessRepo:        clusterAccessRepo,
		tenantClusterBindingRepo: tenantClusterBindingRepo,
	}
}

func NewClusterAccessClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.ClusterAccessClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithInsecure(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewClusterAccessClient(conn), nil
}

func (s *clusterAccessServer) GetClusterAccessByTenantId(id *wrapperspb.StringValue, stream api.ClusterAccess_GetClusterAccessByTenantIdServer) error {
	tenantId, err := uuid.Parse(id.GetValue())
	if err != nil {
		return err
	}

	clusters, err := s.clusterAccessRepo.GetClustersAccessibleByTenantId(stream.Context(), tenantId)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, c := range clusters {
		err := stream.Send(c)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

func (s *clusterAccessServer) GetClusterAccessByUserId(id *wrapperspb.StringValue, stream api.ClusterAccess_GetClusterAccessByUserIdServer) error {
	userId, err := uuid.Parse(id.GetValue())
	if err != nil {
		return err
	}

	clusters, err := s.clusterAccessRepo.GetClustersAccessibleByUserId(stream.Context(), userId)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, c := range clusters {
		err := stream.Send(c)
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

func (s *clusterAccessServer) GetTenantClusterMappingsByTenantId(id *wrapperspb.StringValue, stream api.ClusterAccess_GetTenantClusterMappingsByTenantIdServer) error {
	tenantId, err := uuid.Parse(id.GetValue())
	if err != nil {
		return err
	}

	bindings, err := s.tenantClusterBindingRepo.GetByTenantId(stream.Context(), tenantId)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, t := range bindings {
		err := stream.Send(t.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

func (s *clusterAccessServer) GetTenantClusterMappingsByClusterId(id *wrapperspb.StringValue, stream api.ClusterAccess_GetTenantClusterMappingsByClusterIdServer) error {
	clusterId, err := uuid.Parse(id.GetValue())
	if err != nil {
		return err
	}

	bindings, err := s.tenantClusterBindingRepo.GetByClusterId(stream.Context(), clusterId)
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, t := range bindings {
		err := stream.Send(t.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

func (s *clusterAccessServer) GetTenantClusterMappingByTenantAndClusterId(ctx context.Context, request *api.GetClusterMappingRequest) (*projections.TenantClusterBinding, error) {
	tenantId, err := uuid.Parse(request.GetTenantId())
	if err != nil {
		return nil, err
	}

	clusterId, err := uuid.Parse(request.GetClusterId())
	if err != nil {
		return nil, err
	}

	binding, err := s.tenantClusterBindingRepo.GetByTenantAndClusterId(ctx, tenantId, clusterId)
	if err != nil {
		return nil, err
	}
	return binding.Proto(), nil
}
