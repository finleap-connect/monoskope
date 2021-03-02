package commands

import (
	"context"

	"github.com/google/uuid"

	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateUserCommand is a command for creating a user.
type CreateUserCommand struct {
	aggregateId uuid.UUID
	cmdData.CreateUserCommandData
}

func NewCreateUserCommand() *CreateUserCommand {
	return &CreateUserCommand{
		aggregateId:           uuid.New(),
		CreateUserCommandData: cmdData.CreateUserCommandData{},
	}
}

func (c *CreateUserCommand) AggregateID() uuid.UUID          { return c.aggregateId }
func (c *CreateUserCommand) AggregateType() es.AggregateType { return aggregates.User }
func (c *CreateUserCommand) CommandType() es.CommandType     { return commands.CreateUser }
func (c *CreateUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserCommandData)
}
func (c *CreateUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithSubject(c.GetEmail()),                      // Allows user to create themselves
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows user to create themselves
	}
}
