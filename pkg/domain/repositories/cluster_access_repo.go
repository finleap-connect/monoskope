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

package repositories

type clusterAccessRepository struct {
	clusterRepo         ReadOnlyClusterRepository
	tenantRepo          ReadOnlyTenantRepository
	userRepo            ReadOnlyUserRepository
	userRoleBindingRepo ReadOnlyUserRoleBindingRepository
}

// ReadOnlyClusterAccessRepository is a repository for reading accesses to a cluster.
type ReadOnlyClusterAccessRepository interface {
}

// NewClusterAccessRepository creates a repository for reading cluster access projections.
func NewClusterAccessRepository(tenantRepo ReadOnlyTenantRepository, clusterRepo ReadOnlyClusterRepository, userRepo ReadOnlyUserRepository, userRoleBindingRepo ReadOnlyUserRoleBindingRepository) ReadOnlyClusterAccessRepository {
	return &clusterAccessRepository{
		clusterRepo:         clusterRepo,
		tenantRepo:          tenantRepo,
		userRepo:            userRepo,
		userRoleBindingRepo: userRoleBindingRepo,
	}
}

// func (r *clusterAccessRepository) GetClustersAccessibleByUserId(ctx context.Context, id uuid.UUID) ([]*projections.UserClusterRole, error) {
// 	roleBindings, err := r.userRoleBindingRepo.ByUserId(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// }
