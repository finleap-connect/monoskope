package projections

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type User struct {
	Email string
}

func NewUser() *User {
	return &User{}
}

func (u *User) ID() uuid.UUID {
	return uuid.Nil
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
		domain.User,
		domain.UserRoleBinding,
	}
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, e es.Event, p es.Projection) (es.Projection, error) {
	return p, nil
}
