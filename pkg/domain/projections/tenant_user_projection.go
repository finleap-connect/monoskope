package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type TenantUser struct {
	*DomainProjection
	*projections.TenantUser
}

func NewTenantUserProjection(tenantId uuid.UUID, user *User, rolebinding *UserRoleBinding) *TenantUser {
	return &TenantUser{
		TenantUser: &projections.TenantUser{
			Id:         user.Id,
			Name:       user.Name,
			Email:      user.Email,
			TenantRole: rolebinding.Role,
			TenantId:   tenantId.String(),
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *TenantUser) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *TenantUser) Proto() *projections.TenantUser {
	return p.TenantUser
}
