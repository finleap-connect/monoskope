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
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	domain_projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/google/uuid"
)

type clusterAccessRepository struct {
	clusterRepo              ClusterRepository
	userRoleBindingRepo      UserRoleBindingRepository
	tenantClusterBindingRepo TenantClusterBindingRepository
}

// ClusterAccessRepository is a repository for reading accesses to a cluster.
type ClusterAccessRepository interface {
	// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
	GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) ([]*projections.ClusterAccess, error)
}

// NewClusterAccessRepository creates a repository for reading cluster access projections.
func NewClusterAccessRepository(tenantClusterBindingRepo TenantClusterBindingRepository, clusterRepo ClusterRepository, userRoleBindingRepo UserRoleBindingRepository) ClusterAccessRepository {
	return &clusterAccessRepository{
		clusterRepo:              clusterRepo,
		userRoleBindingRepo:      userRoleBindingRepo,
		tenantClusterBindingRepo: tenantClusterBindingRepo,
	}
}

// getClustersByBindings returns all clusters part of the bindings.
func (r *clusterAccessRepository) getClustersByBindings(ctx context.Context, bindings []*domain_projections.TenantClusterBinding, roles []string) (clusters []*projections.ClusterAccess, err error) {
	for _, clusterBinding := range bindings {
		var cluster *domain_projections.Cluster
		cluster, err = r.clusterRepo.ById(ctx, uuid.MustParse(clusterBinding.ClusterId))
		if err != nil {
			return
		}
		clusters = append(clusters, &projections.ClusterAccess{Cluster: cluster.Cluster, Roles: roles})
	}
	return
}

// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
func (r *clusterAccessRepository) GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) (clusters []*projections.ClusterAccess, err error) {
	// get all rolebindings of the user
	var roleBindings []*domain_projections.UserRoleBinding
	roleBindings, err = r.userRoleBindingRepo.ByUserId(ctx, id)
	if err != nil {
		return
	}

	// check if user is system admin
	var isSystemAdmin = false
	for _, roleBinding := range roleBindings {
		isSystemAdmin = isSystemAdmin || (roleBinding.Scope == string(scopes.System) && roleBinding.Role == string(roles.Admin))
	}

	if isSystemAdmin { // system admins have access to all clusters
		var c []*domain_projections.Cluster
		c, err = r.clusterRepo.AllWith(ctx, false)
		if err != nil {
			return
		}
		for _, cluster := range c {
			clusters = append(clusters,
				&projections.ClusterAccess{Cluster: cluster.Cluster, Roles: []string{
					string(k8s.DefaultRole),
					string(k8s.AdminRole),
				}})
		}
	} else { // regular users have access based on tenant membership
		for _, roleBinding := range roleBindings {
			// search rolebindings for tenant scoped bindings
			if roleBinding.Scope == string(scopes.Tenant) {
				// Set roles within cluster
				var k8sRoles []string = []string{
					string(k8s.DefaultRole),
				}
				// System admins are admins in k8s clusters too
				if isSystemAdmin {
					k8sRoles = append(k8sRoles, string(k8s.AdminRole))
				}
				if roleBinding.Role == string(roles.OnCall) {
					k8sRoles = append(k8sRoles, string(k8s.OnCallRole))
				}

				// get accessible cluster by tenant and append
				var bindings []*domain_projections.TenantClusterBinding
				bindings, err = r.tenantClusterBindingRepo.GetByTenantId(ctx, uuid.MustParse(roleBinding.GetResource()))
				if err != nil {
					return
				}
				clusters, err = r.getClustersByBindings(ctx, bindings, k8sRoles)
			}
		}
	}
	return
}
