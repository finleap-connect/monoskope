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

// clusterServer is the implementation of the ClusterService API
type clusterServer struct {
	api.UnimplementedClusterServer

	repoCluster repositories.ClusterRepository
}

// NewClusterServiceServer returns a new configured instance of clusterServiceServer
func NewClusterServer(clusterRepo repositories.ClusterRepository) *clusterServer {
	return &clusterServer{
		repoCluster: clusterRepo,
	}
}

func NewClusterClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.ClusterClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithInsecure(queryHandlerAddr).
		WithOpenTelemetry().
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewClusterClient(conn), nil
}

// GetById returns the cluster found by the given id.
func (s *clusterServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.Cluster, error) {
	uid, err := uuid.Parse(id.GetValue())
	if err != nil {
		return nil, err
	}

	cluster, err := s.repoCluster.ById(ctx, uid)
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
	clusters, err := s.repoCluster.AllWith(stream.Context(), request.GetIncludeDeleted())
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
