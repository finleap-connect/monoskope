package projections

import (
	"github.com/google/uuid"
)

type User struct {
	id    uuid.UUID
	name  string
	email string
	roles []*UserRoleBinding
}

func (u *User) ID() uuid.UUID {
	return u.id
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
	return &User{id: id, name: name, email: email, roles: roles}
}
