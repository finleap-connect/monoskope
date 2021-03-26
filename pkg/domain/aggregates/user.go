package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	eventData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*es.BaseAggregate
	aggregateManager es.AggregateManager
	Email            string
	Name             string
}

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(id uuid.UUID, aggregateManager es.AggregateManager) *UserAggregate {
	return &UserAggregate{
		BaseAggregate:    es.NewBaseAggregate(aggregates.User, id),
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		return a.createUser(ctx, cmd)
	}
	return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserCreated:
		err := a.userCreated(event)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}

	return nil
}

func containsUser(values []es.Aggregate, emailAddress string) bool {
	for _, value := range values {
		d, ok := value.(*UserAggregate)
		if ok {
			if d.Email == emailAddress {
				return true
			}
		}
	}
	return false
}

// createUser handle the command
func (a *UserAggregate) createUser(ctx context.Context, cmd *commands.CreateUserCommand) error {
	// Get all aggregates of same type
	aggregates, err := a.aggregateManager.All(ctx, a.Type())
	if err != nil {
		return err
	}

	// Check if user already exists
	if !containsUser(aggregates, cmd.GetEmail()) {
		// User does not exist yet, creating...
		eventData := &eventData.UserCreatedEventData{
			Email: cmd.GetEmail(),
			Name:  cmd.GetName(),
		}
		_ = a.AppendEvent(ctx, events.UserCreated, es.ToEventDataFromProto(eventData))
		return nil
	} else {
		return domainErrors.ErrUserAlreadyExists
	}
}

// userCreated handle the event
func (a *UserAggregate) userCreated(event es.Event) error {
	data := &eventData.UserCreatedEventData{}
	if err := event.Data().ToProto(data); err != nil {
		return err
	}

	a.Email = data.Email
	a.Name = data.Name

	return nil
}
