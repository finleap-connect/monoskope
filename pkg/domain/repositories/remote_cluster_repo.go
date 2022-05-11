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

package repositories

import (
	"context"

	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type remoteClusterRepository struct {
	clusterClient api.ClusterClient
}

// NewRemoteClusterRepository creates a repository for reading user projections.
func NewRemoteClusterRepository(clusterClient api.ClusterClient) ReadOnlyClusterRepository {
	return &remoteClusterRepository{
		clusterClient: clusterClient,
	}
}

func (r *remoteClusterRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Cluster, error) {
	panic("not implemented")
}

func (r *remoteClusterRepository) GetBootstrapToken(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *remoteClusterRepository) ByClusterId(ctx context.Context, id string) (*projections.Cluster, error) {
	clusterProto, err := r.clusterClient.GetById(ctx, wrapperspb.String(id))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}
	return &projections.Cluster{Cluster: clusterProto}, nil
}

func (r *remoteClusterRepository) ByClusterName(ctx context.Context, name string) (*projections.Cluster, error) {
	clusterProto, err := r.clusterClient.GetByName(ctx, wrapperspb.String(name))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}
	return &projections.Cluster{Cluster: clusterProto}, nil
}
