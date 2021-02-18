package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*es.BaseAggregate
	Email string
	Name  string
}

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(id uuid.UUID) *UserAggregate {
	return &UserAggregate{
		BaseAggregate: es.NewBaseAggregate(aggregates.User, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		// TODO: check if user already exists
		if ed, err := es.ToEventDataFromProto(&ed.UserCreatedEventData{Email: cmd.GetEmail(), Name: cmd.GetName()}); err != nil {
			return err
		} else if err = a.ApplyEvent(a.AppendEvent(ctx, events.UserCreated, ed)); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserCreated:
		data := &ed.UserCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Email = data.Email
		a.Name = data.Name
	default:
		return fmt.Errorf("couldn't handle event")
	}
	return nil
}
