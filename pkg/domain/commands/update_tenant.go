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
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewUpdateTenantCommand)
}

// UpdateTenantCommand is a command for updating a tenant.
type UpdateTenantCommand struct {
	*es.BaseCommand
	cmdData.UpdateTenantCommandData
}

// NewUpdateTenantCommand creates an UpdateTenantCommand.
func NewUpdateTenantCommand(id uuid.UUID) es.Command {
	return &UpdateTenantCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Tenant, commands.UpdateTenant),
	}
}

func (c *UpdateTenantCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.UpdateTenantCommandData)
}
