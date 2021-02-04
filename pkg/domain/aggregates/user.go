package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	types "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*es.BaseAggregate
	email string
	name  string
}

// AggregateType returns the type of the aggregate.
func (c *UserAggregate) AggregateType() es.AggregateType { return types.UserType }

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(id uuid.UUID) *UserAggregate {
	return &UserAggregate{
		BaseAggregate: es.NewBaseAggregate(types.UserType, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		// TODO: check if user already exists
		// TODO: check if user is allowed to do this
		if ed, err := es.ToEventDataFromProto(&ed.UserCreatedEventData{Email: cmd.UserMetadata.Email, Name: cmd.UserMetadata.Name}); err != nil {
			return err
		} else if err = a.ApplyEvent(a.AppendEvent(types.UserCreatedType, ed)); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case types.UserCreatedType:
		data := &ed.UserCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.email = data.Email
		a.name = data.Name
	}
	return nil
}
