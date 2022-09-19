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
	tenantRepo               TenantRepository
	clusterRepo              ClusterRepository
	userRoleBindingRepo      UserRoleBindingRepository
	tenantClusterBindingRepo TenantClusterBindingRepository
}

// ClusterAccessRepository is a repository for reading accesses to a cluster.
type ClusterAccessRepository interface {
	// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
	GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) ([]*projections.ClusterAccess, error)
	// GetClustersAccessibleByUserIdV2 returns all clusters accessible by a user identified by user id
	GetClustersAccessibleByUserIdV2(ctx context.Context, id uuid.UUID) ([]*projections.ClusterAccessV2, error)
}

// NewClusterAccessRepository creates a repository for reading cluster access projections.
func NewClusterAccessRepository(tenantClusterBindingRepo TenantClusterBindingRepository, clusterRepo ClusterRepository, userRoleBindingRepo UserRoleBindingRepository, tenantRepo TenantRepository) ClusterAccessRepository {
	return &clusterAccessRepository{
		clusterRepo:              clusterRepo,
		userRoleBindingRepo:      userRoleBindingRepo,
		tenantClusterBindingRepo: tenantClusterBindingRepo,
		tenantRepo:               tenantRepo,
	}
}

// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
func (r *clusterAccessRepository) GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) (clusters []*projections.ClusterAccess, err error) {
	clustersV2, clustersV2err := r.GetClustersAccessibleByUserIdV2(ctx, id)
	if clustersV2err != nil {
		return nil, clustersV2err
	}

	for _, clusterV2 := range clustersV2 {
		var roles []string
		for _, clusterRole := range clusterV2.ClusterRoles {
			roles = append(roles, clusterRole.Role)
		}
		clusters = append(clusters, &projections.ClusterAccess{
			Cluster: clusterV2.Cluster,
			Roles:   roles,
		})
	}
	return
}

// GetClustersAccessibleByUserId returns all clusters accessible by a user identified by user id
func (r *clusterAccessRepository) GetClustersAccessibleByUserIdV2(ctx context.Context, id uuid.UUID) (clusters []*projections.ClusterAccessV2, err error) {
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
				&projections.ClusterAccessV2{
					Cluster: cluster.Cluster,
					ClusterRoles: []*projections.ClusterRole{
						{
							Scope: projections.ClusterRole_CLUSTER,
							Role:  string(k8s.DefaultRole),
						},
						{
							Scope: projections.ClusterRole_CLUSTER,
							Role:  string(k8s.AdminRole),
						},
						{
							Scope: projections.ClusterRole_CLUSTER,
							Role:  string(k8s.OnCallRole),
						},
					},
				})
		}
		return
	}

	tenantMap := make(map[string]bool)

	// regular users have access based on tenant membership
	for _, binding := range roleBindings {
		// search rolebindings for tenant scoped bindings
		if binding.Scope == string(scopes.Tenant) {
			// check if we already had this tenant
			if _, ok := tenantMap[binding.Resource]; ok {
				continue
			}
			tenantMap[binding.Resource] = true
			tenantId := uuid.MustParse(binding.Resource)

			var tenant *domain_projections.Tenant
			tenant, err = r.tenantRepo.ById(ctx, tenantId)
			if err != nil {
				return
			}
			// Skip deleted tenants
			if tenant.Metadata.Deleted != nil {
				continue
			}

			var tenantBindings []*domain_projections.UserRoleBinding
			tenantBindings, err = r.userRoleBindingRepo.ByUserIdScopeAndResource(ctx, id, scopes.Tenant, binding.Resource)
			if err != nil {
				return
			}

			// Set roles within cluster
			var k8sRoles = []*projections.ClusterRole{
				{
					Scope: projections.ClusterRole_CLUSTER,
					Role:  string(k8s.DefaultRole),
				},
			}

			for _, tenantBinding := range tenantBindings {
				if tenantBinding.Role == string(roles.Admin) {
					k8sRoles = append(k8sRoles, &projections.ClusterRole{
						Scope: projections.ClusterRole_TENANT,
						Role:  string(k8s.AdminRole),
					})
				}
				if tenantBinding.Role == string(roles.OnCall) {
					k8sRoles = append(k8sRoles, &projections.ClusterRole{
						Scope: projections.ClusterRole_TENANT,
						Role:  string(k8s.OnCallRole),
					})
				}
			}

			// get accessible cluster by tenant and append
			var tenantClusterBindings []*domain_projections.TenantClusterBinding
			tenantClusterBindings, err = r.tenantClusterBindingRepo.GetByTenantId(ctx, tenantId)
			if err != nil {
				return
			}

			for _, tcb := range tenantClusterBindings {
				var cluster *domain_projections.Cluster
				cluster, err = r.clusterRepo.ById(ctx, uuid.MustParse(tcb.ClusterId))
				if err != nil {
					return
				}

				// Skip deleted clusters
				if cluster.Metadata.Deleted != nil {
					continue
				}

				clusters = append(clusters, &projections.ClusterAccessV2{Cluster: cluster.Cluster, ClusterRoles: k8sRoles})
			}
		}
	}

	return
}
