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

package projections

import (
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/google/uuid"
)

type TenantUser struct {
	*DomainProjection
	*projections.TenantUser
}

func NewTenantUserProjection(tenantId uuid.UUID, user *User, rolebindings []*UserRoleBinding) *TenantUser {
	dp := NewDomainProjection()
	tu := &TenantUser{
		DomainProjection: dp,
		TenantUser: &projections.TenantUser{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			TenantId: tenantId.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
	for _, roleBinding := range rolebindings {
		tu.TenantRoles = append(tu.TenantRoles, roleBinding.Role)
	}
	return tu
}

// ID implements the ID method of the Aggregate interface.
func (p *TenantUser) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *TenantUser) Proto() *projections.TenantUser {
	return p.TenantUser
}
