// Copyright 2021 Monoskope Authors
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

	"github.com/golang/protobuf/ptypes/wrappers"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// clusterServer is the implementation of the ClusterService API
type clusterServer struct {
	api.UnimplementedClusterServer

	repoCluster repositories.ReadOnlyClusterRepository
}

// NewClusterServiceServer returns a new configured instance of clusterServiceServer
func NewClusterServer(clusterRepo repositories.ReadOnlyClusterRepository) *clusterServer {
	return &clusterServer{
		repoCluster: clusterRepo,
	}
}

func NewClusterClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.ClusterClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewClusterClient(conn), nil
}

// GetById returns the cluster found by the given id.
func (s *clusterServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.Cluster, error) {
	cluster, err := s.repoCluster.ByClusterId(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}
	return cluster.Proto(), nil
}

// GetByName returns the cluster found by the given name.
func (s *clusterServer) GetByName(ctx context.Context, name *wrappers.StringValue) (*projections.Cluster, error) {
	cluster, err := s.repoCluster.ByClusterName(ctx, name.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return cluster.Proto(), nil
}

// GetAll returns all clusters.
func (s *clusterServer) GetAll(request *api.GetAllRequest, stream api.Cluster_GetAllServer) error {
	clusters, err := s.repoCluster.GetAll(stream.Context(), request.GetIncludeDeleted())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, c := range clusters {
		err := stream.Send(c.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

// GetBootstrapToken returns the bootstrap token for the cluster with the given id.
func (s *clusterServer) GetBootstrapToken(ctx context.Context, id *wrappers.StringValue) (*wrappers.StringValue, error) {
	token, err := s.repoCluster.GetBootstrapToken(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{Value: token}, nil
}
