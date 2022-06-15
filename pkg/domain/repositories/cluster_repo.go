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

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

type clusterRepository struct {
	DomainRepository[*projections.Cluster]
}

// ClusterRepository is a repository for reading and writing cluster projections.
type ClusterRepository interface {
	DomainRepository[*projections.Cluster]
	// ByName searches for the a tenant projection by it's name
	ByClusterName(context.Context, string) (*projections.Cluster, error)
	// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID
	GetBootstrapToken(context.Context, string) (string, error)
}

// NewClusterRepository creates a repository for reading and writing cluster projections.
func NewClusterRepository(repository es.Repository[*projections.Cluster]) ClusterRepository {
	return &clusterRepository{
		NewDomainRepository(repository),
	}
}

// ByClusterName searches for a cluster projection by its name.
func (r *clusterRepository) ByClusterName(ctx context.Context, clusterName string) (*projections.Cluster, error) {
	ps, err := r.AllWith(ctx, true)
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

// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID.
func (r *clusterRepository) GetBootstrapToken(ctx context.Context, id string) (string, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	cluster, err := r.ById(ctx, uuid)
	if err != nil {
		return "", err
	}
	return cluster.BootstrapToken, nil
}
