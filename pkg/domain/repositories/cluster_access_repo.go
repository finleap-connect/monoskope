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

	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/google/uuid"
)

type clusterAccessRepository struct {
	clusterRepo              ReadOnlyClusterRepository
	userRoleBindingRepo      ReadOnlyUserRoleBindingRepository
	tenantClusterBindingRepo ReadOnlyTenantClusterBindingRepository
}

// ReadOnlyClusterAccessRepository is a repository for reading accesses to a cluster.
type ReadOnlyClusterAccessRepository interface {
	// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
	GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) ([]*projections.Cluster, error)
	// GetClustersAccessibleByTenantId returns all clusters accessible by a tenant identified by tenant id
	GetClustersAccessibleByTenantId(ctx context.Context, id uuid.UUID) ([]*projections.Cluster, error)
}

// NewClusterAccessRepository creates a repository for reading cluster access projections.
func NewClusterAccessRepository(tenantClusterBindingRepo ReadOnlyTenantClusterBindingRepository, clusterRepo ReadOnlyClusterRepository, userRoleBindingRepo ReadOnlyUserRoleBindingRepository) ReadOnlyClusterAccessRepository {
	return &clusterAccessRepository{
		clusterRepo:              clusterRepo,
		userRoleBindingRepo:      userRoleBindingRepo,
		tenantClusterBindingRepo: tenantClusterBindingRepo,
	}
}

func (r *clusterAccessRepository) GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) ([]*projections.Cluster, error) {
	roleBindings, err := r.userRoleBindingRepo.ByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	var clusters []*projections.Cluster
	for _, roleBinding := range roleBindings {
		if roleBinding.Scope == scopes.Tenant.String() {
			tenantClusterBinding, err := r.tenantClusterBindingRepo.GetByTenantId(ctx, uuid.MustParse(roleBinding.GetResource()))
			if err != nil {
				return nil, err
			}

			for _, clusterBinding := range tenantClusterBinding {
				cluster, err := r.clusterRepo.ByClusterId(ctx, clusterBinding.ClusterId)
				if err != nil {
					return nil, err
				}
				clusters = append(clusters, cluster.Cluster)
			}
		}
	}
	return clusters, nil
}

func (r *clusterAccessRepository) GetClustersAccessibleByTenantId(ctx context.Context, id uuid.UUID) ([]*projections.Cluster, error) {
	bindings, err := r.tenantClusterBindingRepo.GetByTenantId(ctx, id)
	if err != nil {
		return nil, err
	}

	var clusters []*projections.Cluster
	for _, clusterBinding := range bindings {
		cluster, err := r.clusterRepo.ByClusterId(ctx, clusterBinding.ClusterId)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, cluster.Cluster)
	}
	return clusters, nil
}
