package commands

import (
	"context"

	"github.com/google/uuid"

	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateUserCommand is a command for creating a user.
type CreateUserCommand struct {
	aggregateId uuid.UUID
	cmd.CreateUserCommand
}

func (c *CreateUserCommand) AggregateID() uuid.UUID          { return c.aggregateId }
func (c *CreateUserCommand) AggregateType() es.AggregateType { return aggregates.User }
func (c *CreateUserCommand) CommandType() es.CommandType     { return commands.CreateUser }
func (c *CreateUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserCommand)
}
func (c *CreateUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		{Subject: c.GetUserMetadata().GetEmail()}, // Allows user to create themselfes
	}
}
