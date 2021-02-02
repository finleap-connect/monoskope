package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	UserType        es.AggregateType = "User"
	CreateUserType  es.CommandType   = "CreateUser"
	UserCreatedType es.EventType     = "UserCreated"
)

// CreateUserCommand is a command for creating a user.
type CreateUserCommand struct {
	aggregateId uuid.UUID
	cmd.CreateUserCommand
}

func (c *CreateUserCommand) AggregateID() uuid.UUID          { return c.aggregateId }
func (c *CreateUserCommand) AggregateType() es.AggregateType { return UserType }
func (c *CreateUserCommand) CommandType() es.CommandType     { return CreateUserType }
func (c *CreateUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserCommand)
}

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*es.AggregateBase
	email string
	name  string
}

func (c *UserAggregate) AggregateType() es.AggregateType { return UserType }

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate() *UserAggregate {
	return &UserAggregate{
		AggregateBase: es.NewAggregateBase(UserType, uuid.New()),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *CreateUserCommand:
		// TODO: check if user already exists
		// TODO: check if user is allowed to do this
		if ed, err := es.ToEventDataFromProto(&ed.UserCreatedEventData{Email: cmd.UserMetadata.Email, Name: cmd.UserMetadata.Name}); err != nil {
			return err
		} else if err = a.ApplyEvent(a.AppendEvent(UserCreatedType, ed)); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case UserCreatedType:
		data := &ed.UserCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.email = data.Email
		a.name = data.Name
	}
	return nil
}

type User struct {
	Name  string
	Email string
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
		UserType,
		UserRoleBindingType,
	}
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, e es.Event, p es.Projection) (es.Projection, error) {
	return p, nil
}
