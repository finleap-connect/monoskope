package projections

import (
	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type UserRoleBinding struct {
	id       uuid.UUID
	role     es.Role
	scope    es.Scope
	resource string
}

func (u *UserRoleBinding) ID() uuid.UUID {
	return u.id
}

func (u *UserRoleBinding) Role() es.Role {
	return u.role
}

func (u *UserRoleBinding) Scope() es.Scope {
	return u.scope
}

func (u *UserRoleBinding) Resource() string {
	return u.resource
}

func NewUserRoleBinding(id uuid.UUID, role es.Role, scope es.Scope, resource string) *UserRoleBinding {
	return &UserRoleBinding{id: id, role: role, scope: scope, resource: resource}
}
