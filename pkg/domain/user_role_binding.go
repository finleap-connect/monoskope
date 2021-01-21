package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	apicmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

const (
	AddRoleToUser   evs.CommandType   = "AddRoleToUser"
	UserRoleBinding evs.AggregateType = "UserRoleBinding"
	UserRoleAdded   evs.EventType     = "UserRoleAdded"
)

func init() {
	if err := evs.Registry.RegisterCommand(addRoleToUserCommandFactory); err != nil {
		panic(err)
	}
}

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	AggID uuid.UUID
	apicmd.AddRoleToUser
}

var addRoleToUserCommandFactory = func() evs.Command { return &AddRoleToUserCommand{} }

func (c *AddRoleToUserCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *AddRoleToUserCommand) AggregateType() evs.AggregateType {
	return UserRoleBinding
}
func (c *AddRoleToUserCommand) CommandType() evs.CommandType { return AddRoleToUser }

// UserRoleBindingAggregate is an aggregate for .
type UserRoleBindingAggregate struct {
	*evs.AggregateBase
	userId  uuid.UUID
	role    string
	context string
}

func (c *UserRoleBindingAggregate) AggregateType() evs.AggregateType {
	return UserRoleBinding
}

func NewUserRoleBindingAggregate(id uuid.UUID) *UserRoleBindingAggregate {
	return &UserRoleBindingAggregate{
		AggregateBase: evs.NewAggregateBase(UserRoleBinding, id),
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd evs.Command) ([]evs.Event, error) {
	var resultingEvents []evs.Event

	switch cmd := cmd.(type) {
	case *AddRoleToUserCommand:
		if event, err := a.HandleAddRoleToUserCommand(cmd); err == nil {
			resultingEvents = append(resultingEvents, event)
		} else {
			return nil, err
		}
		return resultingEvents, nil
	}

	return nil, fmt.Errorf("couldn't handle command")
}

func (a *UserRoleBindingAggregate) HandleAddRoleToUserCommand(cmd *AddRoleToUserCommand) (evs.Event, error) {
	// TODO: Check if user has the right to do this.
	d := &UserRoleAddedData{
		Role:    cmd.GetRole(),
		Context: cmd.GetContext(),
	}
	ed, err := evs.ToEventData(d)
	if err != nil {
		return nil, err
	}
	return a.AppendEvent(UserRoleAdded, ed), nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(ctx context.Context, event evs.Event) error {
	switch event.EventType() {
	case UserRoleAdded:
		// TODO: build up aggregate from events
		a.userId = uuid.Nil
		a.role = ""
		a.context = ""
	}
	return nil
}

type UserRoleAddedData struct {
	Role    string `json:",omitempty"`
	Context string `json:",omitempty"`
}

func UserRoleAddedDataFromEventData(ed evs.EventData) (*UserRoleAddedData, error) {
	data := &UserRoleAddedData{}
	if err := json.Unmarshal(ed, data); err != nil {
		return nil, err
	}
	return data, nil
}
