package domain

import (
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
	ID uuid.UUID
	evs.AggregateBase
}

func (c *UserRoleBindingAggregate) AggregateType() evs.AggregateType {
	return UserRoleBinding
}

// func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd evs.Command) error {
// 	switch cmd := cmd.(type) {
// 	case *AddRoleToUserCommand:
// 		//TODO: Check if that is allowed

// 		// If yes append event
// 		a.AppendEvent(UserRoleAdded,
// 			&UserRoleAddedData{
// 				Role:    cmd.GetRole(),
// 				Context: cmd.GetContext(),
// 			},
// 			time.Now(),
// 		)
// 		return nil
// 	}
// 	return fmt.Errorf("couldn't handle command")
// }

type UserRoleAddedData struct {
	Role    string `json:",omitempty"`
	Context string `json:",omitempty"`
}
