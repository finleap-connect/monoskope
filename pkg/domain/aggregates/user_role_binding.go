package aggregates

import (
	"fmt"

	"github.com/google/uuid"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// UserRoleBindingAggregate is an aggregate for .
type UserRoleBindingAggregate struct {
	*AggregateBase
	userId  uuid.UUID
	role    string
	context string
}

func (c *UserRoleBindingAggregate) AggregateType() AggregateType { return domain.UserRoleBinding }

func NewUserRoleBindingAggregate(id uuid.UUID) *UserRoleBindingAggregate {
	return &UserRoleBindingAggregate{
		AggregateBase: NewAggregateBase(domain.UserRoleBinding, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(cmd Command) ([]Event, error) {
	var resultingEvents []Event

	switch cmd := cmd.(type) {
	case *commands.AddRoleToUserCommand:
		if event, err := a.handleAddRoleToUserCommand(cmd); err == nil {
			resultingEvents = append(resultingEvents, event)
		} else {
			return nil, err
		}
		return resultingEvents, nil
	}

	return nil, fmt.Errorf("couldn't handle command")
}

func (a *UserRoleBindingAggregate) handleAddRoleToUserCommand(cmd *commands.AddRoleToUserCommand) (Event, error) {
	// TODO: Check if user has the right to do this.
	userId, err := uuid.Parse(cmd.GetUserId())
	if err != nil {
		return nil, err
	}

	ed, err := ToEventDataFromProto(&api.UserRoleAddedEventData{
		UserId:  userId.String(),
		Role:    cmd.GetRole(),
		Context: cmd.GetContext(),
	})
	if err != nil {
		return nil, err
	}

	return a.AppendEvent(domain.UserRoleAdded, ed), nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event Event) error {
	switch event.EventType() {
	case domain.UserRoleAdded:
		err := a.applyUserRoleAddedEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *UserRoleBindingAggregate) applyUserRoleAddedEvent(event Event) error {
	data := &api.UserRoleAddedEventData{}
	err := event.Data().ToProto(data)
	if err != nil {
		return err
	}

	userId, err := uuid.Parse(data.UserId)
	if err != nil {
		return err
	}

	a.userId = userId
	a.role = data.Role
	a.context = data.Context

	return nil
}
