package projections

import (
	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type UserRoleBinding struct {
	es.BaseProjection
	id       uuid.UUID
	userId   uuid.UUID
	role     es.Role
	scope    es.Scope
	resource string
}

func (u *UserRoleBinding) ID() uuid.UUID {
	return u.id
}

func (u *UserRoleBinding) UserID() uuid.UUID {
	return u.userId
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

func NewUserRoleBinding(id, userId uuid.UUID, role es.Role, scope es.Scope, resource string) *UserRoleBinding {
	return &UserRoleBinding{BaseProjection: es.NewBaseProjection(id), userId: userId, role: role, scope: scope, resource: resource}
}
