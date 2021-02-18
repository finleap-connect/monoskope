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

// UserRoleBindingAggregate is an aggregate for UserRoleBindings.
type UserRoleBindingAggregate struct {
	*es.BaseAggregate
	userId   uuid.UUID // User to add a role to
	role     es.Role   // Role to add to the user
	scope    es.Scope  // Scope of the role binding
	resource string    // Resource of the role binding
}

// NewUserRoleBindingAggregate creates a new UserRoleBindingAggregate
func NewUserRoleBindingAggregate(id uuid.UUID) *UserRoleBindingAggregate {
	return &UserRoleBindingAggregate{
		BaseAggregate: es.NewBaseAggregate(aggregates.UserRoleBinding, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserRoleBindingCommand:
		return a.handleAddRoleToUserCommand(ctx, cmd)
	}
	return fmt.Errorf("couldn't handle command")
}

// handleAddRoleToUserCommand handles the command
func (a *UserRoleBindingAggregate) handleAddRoleToUserCommand(ctx context.Context, cmd *commands.CreateUserRoleBindingCommand) error {
	ed, err := es.ToEventDataFromProto(&ed.UserRoleAddedEventData{
		UserId:   cmd.GetUserId(),
		Role:     cmd.GetRole(),
		Scope:    cmd.GetScope(),
		Resource: cmd.GetResource(),
	})
	if err != nil {
		return err
	}

	_ = a.AppendEvent(ctx, events.UserRoleBindingCreated, ed)

	return nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserRoleBindingCreated:
		err := a.applyUserRoleAddedEvent(event)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("couldn't handle event")
	}

	return nil
}

// applyUserRoleAddedEvent applies the event on the aggregate
func (a *UserRoleBindingAggregate) applyUserRoleAddedEvent(event es.Event) error {
	data := &ed.UserRoleAddedEventData{}
	err := event.Data().ToProto(data)
	if err != nil {
		return err
	}

	userId, err := uuid.Parse(data.UserId)
	if err != nil {
		return err
	}

	a.userId = userId
	a.role = es.Role(data.Role)
	a.scope = es.Scope(data.Scope)
	a.resource = data.Resource

	return nil
}
