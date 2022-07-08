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

package aggregates

import (
	"context"
	"fmt"
	"strings"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	aggregates "github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

// TenantAggregate is an aggregate for Tenants.
type TenantAggregate struct {
	*DomainAggregateBase
	aggregateManager es.AggregateStore
	name             string
	prefix           string
}

// NewTenantAggregate creates a new TenantAggregate
func NewTenantAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &TenantAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Tenant),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *TenantAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.validate(ctx, cmd); err != nil {
		return nil, err
	}
	return a.execute(ctx, cmd)
}

// validate validates the current state of the aggregate and if a specific command is valid in the current state
func (a *TenantAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateTenantCommand:
		if a.Exists() {
			return domainErrors.ErrTenantAlreadyExists
		}

		// Get all aggregates of same type
		aggregates, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}

		if containsTenant(aggregates, cmd.GetName()) {
			return domainErrors.ErrTenantAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func containsTenant(values []es.Aggregate, name string) bool {
	for _, value := range values {
		d, ok := value.(*TenantAggregate)
		if ok {
			if !d.Deleted() && strings.ToLower(strings.TrimSpace(d.name)) == strings.ToLower(strings.TrimSpace(name)) {
				return true
			}
		}
	}
	return false
}

// execute executes the command after it has successfully been validated
func (a *TenantAggregate) execute(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	switch cmd := cmd.(type) {
	case *commands.CreateTenantCommand:
		ed := es.ToEventDataFromProto(&eventdata.TenantCreated{
			Name:   cmd.GetName(),
			Prefix: cmd.GetPrefix()})
		_ = a.AppendEvent(ctx, events.TenantCreated, ed)
	case *commands.UpdateTenantCommand:
		ed := es.ToEventDataFromProto(&eventdata.TenantUpdated{
			Name: cmd.GetName(),
		})
		_ = a.AppendEvent(ctx, events.TenantUpdated, ed)
	case *commands.DeleteTenantCommand:
		_ = a.AppendEvent(ctx, events.TenantDeleted, nil)
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	return a.DefaultReply(), nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *TenantAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.TenantCreated:
		data := &eventdata.TenantCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.name = data.Name
		a.prefix = data.Prefix
	case events.TenantUpdated:
		data := &eventdata.TenantUpdated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.name = data.GetName().GetValue()
	case events.TenantDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
