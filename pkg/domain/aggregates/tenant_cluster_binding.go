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

package aggregates

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	aggregates "github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// TenantClusterBindingAggregate is an aggregate for TenantClusterBindings.
type TenantClusterBindingAggregate struct {
	*DomainAggregateBase
	aggregateManager es.AggregateStore
	tenantID         uuid.UUID // ID of the referenced tenant
	clusterID        uuid.UUID // ID of the referenced cluster
}

// NewTenantClusterBindingAggregate creates a new TenantClusterBindingAggregate
func NewTenantClusterBindingAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &TenantClusterBindingAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.TenantClusterBinding),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *TenantClusterBindingAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.Authorize(ctx, cmd, uuid.Nil); err != nil {
		return nil, err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return nil, err
	}

	switch cmd := cmd.(type) {
	case *commands.CreateTenantClusterBindingCommand:
		ed := es.ToEventDataFromProto(&eventdata.TenantClusterBindingCreated{
			TenantId:  cmd.GetTenantId(),
			ClusterId: cmd.GetClusterId(),
		})
		_ = a.AppendEvent(ctx, events.TenantClusterBindingCreated, ed)
		reply := &es.CommandReply{
			Id:      a.ID(),
			Version: a.Version(),
		}
		return reply, nil
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// validate validates the current state of the aggregate and if a specific command is valid in the current state
func (a *TenantClusterBindingAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateTenantClusterBindingCommand:
		if a.Exists() {
			return domainErrors.ErrTenantClusterBindingAlreadyExists
		}

		tenantId, err := uuid.Parse(cmd.GetTenantId())
		if err != nil {
			return domainErrors.ErrInvalidArgument("tenant id is invalid")
		}
		clusterId, err := uuid.Parse(cmd.GetClusterId())
		if err != nil {
			return domainErrors.ErrInvalidArgument("cluster id is invalid")
		}

		tenantAggregate, err := a.aggregateManager.Get(ctx, aggregates.Tenant, tenantId)
		if err != nil {
			return err
		}
		if !tenantAggregate.Exists() || tenantAggregate.Deleted() {
			return domainErrors.ErrTenantNotFound
		}

		clusterAggregate, err := a.aggregateManager.Get(ctx, aggregates.Cluster, clusterId)
		if err != nil {
			return err
		}
		if !clusterAggregate.Exists() || clusterAggregate.Deleted() {
			return domainErrors.ErrClusterNotFound
		}

		// Get all aggregates of same type
		aggs, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}
		if containsTenantClusterBinding(aggs, tenantId, clusterId) {
			return domainErrors.ErrTenantClusterBindingAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func containsTenantClusterBinding(values []es.Aggregate, tenantId, clusterId uuid.UUID) bool {
	for _, value := range values {
		d, ok := value.(*TenantClusterBindingAggregate)
		if ok {
			if !d.Deleted() && d.tenantID == tenantId && d.clusterID == clusterId {
				return true
			}
		}
	}
	return false
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *TenantClusterBindingAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.TenantClusterBindingCreated:
		ed := new(eventdata.TenantClusterBindingCreated)
		err := event.Data().ToProto(ed)
		if err != nil {
			return err
		}

		tenantId, err := uuid.Parse(ed.GetTenantId())
		if err != nil {
			return err
		}
		clusterId, err := uuid.Parse(ed.GetClusterId())
		if err != nil {
			return err
		}

		a.tenantID = tenantId
		a.clusterID = clusterId
	case events.TenantClusterBindingDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
