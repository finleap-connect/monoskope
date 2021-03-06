// Copyright 2022 Monoskope Authors
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
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewDeleteUserRoleBindingCommand)
}

// DeleteUserRoleBindingCommand is a command for removing a role from a user.
type DeleteUserRoleBindingCommand struct {
	*es.BaseCommand
}

func NewDeleteUserRoleBindingCommand(id uuid.UUID) es.Command {
	return &DeleteUserRoleBindingCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.UserRoleBinding, commands.DeleteUserRoleBinding),
	}
}
func (c *DeleteUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return nil
}
