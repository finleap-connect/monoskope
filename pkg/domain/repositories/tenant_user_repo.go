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

	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
)

type tenantuserRepository struct {
	userRepo            UserRepository
	userRoleBindingRepo UserRoleBindingRepository
	tenantRepo          TenantRepository
}

// TenantUserRepository is a repository for reading users of a tenant.
type TenantUserRepository interface {
	// GetTenantUsersById searches for users belonging to a tenant.
	GetTenantUsersById(context.Context, uuid.UUID) ([]*projections.TenantUser, error)
}

// NewTenantUserRepository creates a repository for reading and writing tenantuser projections.
func NewTenantUserRepository(userRepo UserRepository, userRoleBindingRepo UserRoleBindingRepository, tenantRepo TenantRepository) TenantUserRepository {
	return &tenantuserRepository{
		userRepo:            userRepo,
		userRoleBindingRepo: userRoleBindingRepo,
		tenantRepo:          tenantRepo,
	}
}

// GetTenantUsersById searches for users belonging to a tenant.
func (r *tenantuserRepository) GetTenantUsersById(ctx context.Context, id uuid.UUID) ([]*projections.TenantUser, error) {
	roleBindings, err := r.userRoleBindingRepo.ByScopeAndResource(ctx, scopes.Tenant, id)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]bool)
	var tenantUsers []*projections.TenantUser

	for _, binding := range roleBindings {
		user, err := r.userRepo.ByUserId(ctx, uuid.MustParse(binding.UserId))
		if err != nil {
			return nil, err
		}
		// skip deleted users
		if user.IsDeleted() {
			continue
		}

		// Skip deleted tenants
		tenant, err := r.tenantRepo.ById(ctx, uuid.MustParse(binding.Resource))
		if err != nil {
			return nil, err
		}
		if tenant.IsDeleted() {
			continue
		}

		// check if we already had this user for this tenant
		if _, ok := userMap[binding.UserId]; !ok {
			bindings, err := r.userRoleBindingRepo.ByUserIdScopeAndResource(ctx, user.ID(), scopes.Tenant, binding.Resource)
			if err != nil {
				return nil, err
			}

			userMap[user.Id] = true
			tenantUsers = append(tenantUsers, projections.NewTenantUserProjection(id, user, bindings))
		}
	}

	return tenantUsers, nil
}
