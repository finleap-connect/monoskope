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

	"github.com/google/uuid"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewCreateClusterCommand)
}

// CreateClusterCommand is a command for creating a cluster.
type CreateClusterCommand struct {
	*es.BaseCommand
	cmdData.CreateCluster
}

// NewCreateClusterCommand creates a CreateClusterCommand.
func NewCreateClusterCommand(id uuid.UUID) es.Command {
	return &CreateClusterCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Cluster, commands.CreateCluster),
	}
}

func (c *CreateClusterCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateCluster)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *CreateClusterCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins
	}
}
