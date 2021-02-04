package projections

import (
	"context"

	"github.com/google/uuid"
	types "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
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

type userProjector struct{}

func NewUserProjector() es.Projector {
	return &userProjector{}
}

// EvenTypes returns the EvenTypes for which events should be projected.
func (u *userProjector) EvenTypes() []es.EventType {
	return []es.EventType{}
}

// AggregateTypes returns the AggregateTypes for which events should be projected.
func (u *userProjector) AggregateTypes() []es.AggregateType {
	return []es.AggregateType{
		types.UserType,
		types.UserRoleBindingType,
	}
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, e es.Event, p es.Projection) (es.Projection, error) {
	return p, nil
}
