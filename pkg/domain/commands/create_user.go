package commands

import (
	"context"

	"github.com/google/uuid"

	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"

	types "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateUserCommand is a command for creating a user.
type CreateUserCommand struct {
	aggregateId uuid.UUID
	cmd.CreateUserCommand
}

func (c *CreateUserCommand) AggregateID() uuid.UUID          { return c.aggregateId }
func (c *CreateUserCommand) AggregateType() es.AggregateType { return types.UserType }
func (c *CreateUserCommand) CommandType() es.CommandType     { return types.CreateUserType }
func (c *CreateUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserCommand)
}
func (c *CreateUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		{Subject: c.GetUserMetadata().GetEmail()}, // Allows user to create themselfes
	}
}
