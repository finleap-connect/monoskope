package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	eventData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*es.BaseAggregate
	aggregateManager es.AggregateManager
	Email            string
	Name             string
	log              logger.Logger
}

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(id uuid.UUID, aggregateManager es.AggregateManager) *UserAggregate {
	return &UserAggregate{
		BaseAggregate:    es.NewBaseAggregate(aggregates.User, id),
		aggregateManager: aggregateManager,
		log:              logger.WithName("user-aggregate"),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		// Get all aggregates of same type
		aggregates, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}

		// Check if user already exists
		if !containsUser(aggregates, cmd.GetEmail()) {
			// User does not exist yet, creating...
			ed, err := es.ToEventDataFromProto(&eventData.UserCreatedEventData{Email: cmd.GetEmail(), Name: cmd.GetName()})
			if err != nil {
				return err
			}
			_ = a.AppendEvent(ctx, events.UserCreated, ed)

			ed, err = es.ToEventDataFromProto(&eventData.UserRoleAddedEventData{
				UserId: cmd.AggregateID().String(),
				Role:   roles.User.String(),
				Scope:  scopes.System.String(),
			})
			if err != nil {
				return err
			}
			_ = a.AppendEvent(ctx, events.UserRoleBindingCreated, ed)
			return nil
		} else {
			return domainErrors.ErrUserAlreadyExists
		}
	}
	return fmt.Errorf("couldn't handle command")
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserCreated:
		data := &eventData.UserCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Email = data.Email
		a.Name = data.Name
	default:
		a.log.Info("Can not handle event. Ignoring...", "EventType", event.EventType())
		return nil
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
