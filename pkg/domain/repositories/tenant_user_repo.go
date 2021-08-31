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

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
)

type tenantuserRepository struct {
	userRepo            ReadOnlyUserRepository
	userRoleBindingRepo ReadOnlyUserRoleBindingRepository
}

// ReadOnlyTenantUserRepository is a repository for reading users of a tenant.
type ReadOnlyTenantUserRepository interface {
	// GetTenantUsersById searches for users belonging to a tenant.
	GetTenantUsersById(context.Context, uuid.UUID) ([]*projections.TenantUser, error)
}

// NewTenantUserRepository creates a repository for reading and writing tenantuser projections.
func NewTenantUserRepository(userRepo ReadOnlyUserRepository, userRoleBindingRepo ReadOnlyUserRoleBindingRepository) ReadOnlyTenantUserRepository {
	return &tenantuserRepository{
		userRepo:            userRepo,
		userRoleBindingRepo: userRoleBindingRepo,
	}
}

// GetTenantUsersById searches for users belonging to a tenant.
func (r *tenantuserRepository) GetTenantUsersById(ctx context.Context, id uuid.UUID) ([]*projections.TenantUser, error) {
	roleBindings, err := r.userRoleBindingRepo.ByScopeAndResource(ctx, scopes.Tenant, id)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*projections.TenantUser)
	var tenantUsers []*projections.TenantUser

	for _, binding := range roleBindings {
		user, err := r.userRepo.ByUserId(ctx, uuid.MustParse(binding.UserId))
		if err != nil {
			return nil, err
		}

		tu := projections.NewTenantUserProjection(id, user, binding)
		if u, ok := userMap[user.Id]; ok {
			u.TenantRole = fmt.Sprintf("%s,%s", tu.TenantRole, u.TenantRole)
		} else {
			userMap[user.Id] = tu
			tenantUsers = append(tenantUsers, tu)
		}
	}

	return tenantUsers, nil
}
