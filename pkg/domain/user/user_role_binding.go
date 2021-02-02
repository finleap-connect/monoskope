package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	UserRoleBindingType        AggregateType = "UserRoleBinding"
	CreateUserRoleBindingType  CommandType   = "CreateUserRoleBinding"
	UserRoleBindingCreatedType EventType     = "UserRoleBindingCreated"
)

// AddRoleToUser is a command for adding a role to a user.
type CreateUserRoleBindingCommand struct {
	aggregateId uuid.UUID
	cmd.AddRoleToUserCommand
}

func (c *CreateUserRoleBindingCommand) AggregateID() uuid.UUID       { return c.aggregateId }
func (c *CreateUserRoleBindingCommand) AggregateType() AggregateType { return UserRoleBindingType }
func (c *CreateUserRoleBindingCommand) CommandType() CommandType     { return CreateUserRoleBindingType }
func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.AddRoleToUserCommand)
}

// UserRoleBindingAggregate is an aggregate for UserRoleBindings.
type UserRoleBindingAggregate struct {
	*AggregateBase
	userId  uuid.UUID // User to add a role to
	role    string    // Role to add to the user
	context string    // Context of the role binding
}

// AggregateType returns the type of the aggregate.
func (c *UserRoleBindingAggregate) AggregateType() AggregateType { return UserRoleBindingType }

// NewUserRoleBindingAggregate creates a new UserRoleBindingAggregate
func NewUserRoleBindingAggregate(id uuid.UUID) *UserRoleBindingAggregate {
	return &UserRoleBindingAggregate{
		AggregateBase: NewAggregateBase(UserRoleBindingType, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd Command) error {
	switch cmd := cmd.(type) {
	case *CreateUserRoleBindingCommand:
		return a.handleAddRoleToUserCommand(ctx, cmd)
	}
	return fmt.Errorf("couldn't handle command")
}

// handleAddRoleToUserCommand handles the command
func (a *UserRoleBindingAggregate) handleAddRoleToUserCommand(ctx context.Context, cmd *CreateUserRoleBindingCommand) error {
	// TODO: Check if user has the right to do this.
	_, err := metadata.NewDomainMetadataManager(ctx).GetUserInformation() // user issued the command at gateway
	if err != nil {
		return err
	}

	ed, err := ToEventDataFromProto(&ed.UserRoleAddedEventData{
		UserId:  cmd.GetUserId(),
		Role:    cmd.GetRole(),
		Context: cmd.GetContext(),
	})
	if err != nil {
		return err
	}

	_ = a.AppendEvent(UserRoleBindingCreatedType, ed)

	return nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event Event) error {
	switch event.EventType() {
	case UserRoleBindingCreatedType:
		err := a.applyUserRoleAddedEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

// applyUserRoleAddedEvent applies the event on the aggregate
func (a *UserRoleBindingAggregate) applyUserRoleAddedEvent(event Event) error {
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
	a.role = data.Role
	a.context = data.Context

	return nil
}
