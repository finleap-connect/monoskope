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
	"sort"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type clusterRepository struct {
	es.Repository
}

// ClusterRepository is a repository for reading and writing cluster projections.
type ClusterRepository interface {
	es.Repository
	ReadOnlyClusterRepository
}

// ReadOnlyClusterRepository is a repository for reading cluster projections.
type ReadOnlyClusterRepository interface {
	// ById searches for the a tenant projection by it's id.
	ByClusterId(context.Context, string) (*projections.Cluster, error)
	// ByName searches for the a tenant projection by it's name
	ByClusterName(context.Context, string) (*projections.Cluster, error)
	// GetAll searches for all known clusters.
	GetAll(context.Context, bool) ([]*projections.Cluster, error)
	// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID
	GetBootstrapToken(context.Context, string) (string, error)
}

// NewClusterRepository creates a repository for reading and writing cluster projections.
func NewClusterRepository(repository es.Repository) ClusterRepository {
	return &clusterRepository{
		Repository: repository,
	}
}

// ById searches for a cluster by its id.
func (r *clusterRepository) ByClusterId(ctx context.Context, id string) (*projections.Cluster, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	cluster, ok := projection.(*projections.Cluster)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	return cluster, nil
}

// ByClusterName searches for a cluster projection by its name.
func (r *clusterRepository) ByClusterName(ctx context.Context, clusterName string) (*projections.Cluster, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, c := range ps {
		if clusterName == c.Name {
			return c, nil
		}
	}

	return nil, errors.ErrClusterNotFound
}

// GetAll searches for all cluster projections.
func (r *clusterRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Cluster, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	var clusters []*projections.Cluster
	for _, p := range ps {
		if c, ok := p.(*projections.Cluster); ok {
			if !c.GetDeleted().IsValid() || includeDeleted {
				clusters = append(clusters, c)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name > clusters[j].Name
	})
	return clusters, nil
}

// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID.
func (r *clusterRepository) GetBootstrapToken(ctx context.Context, id string) (string, error) {
	cluster, err := r.ByClusterId(ctx, id)
	if err != nil {
		return "", err
	}
	return cluster.BootstrapToken, nil
}
