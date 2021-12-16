// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewDeleteUserCommand)
}

// DeleteUserCommand is a command for deleting a user.
type DeleteUserCommand struct {
	*es.BaseCommand
}

// NewDeleteUserCommand creates a DeleteUserCommand.
func NewDeleteUserCommand(id uuid.UUID) es.Command {
	return &DeleteUserCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.User, commands.DeleteUser),
	}
}

func (c *DeleteUserCommand) SetData(a *anypb.Any) error {
	return nil
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *DeleteUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(es.Role(common.Role_admin.String())).WithScope(es.Scope(common.Scope_system.String())), // Allows system admins
	}
}
