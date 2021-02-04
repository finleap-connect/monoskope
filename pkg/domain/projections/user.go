package projections

import (
	"github.com/google/uuid"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type User struct {
	es.BaseProjection
	name  string
	email string
	roles []*UserRoleBinding
}

func (u *User) Name() string {
	return u.name
}
func (u *User) Email() string {
	return u.email
}
func (u *User) Roles() []*UserRoleBinding {
	return u.roles
}

func NewUser(id uuid.UUID, name, email string, roles []*UserRoleBinding) *User {
	return &User{BaseProjection: es.NewBaseProjection(id), name: name, email: email, roles: roles}
}
