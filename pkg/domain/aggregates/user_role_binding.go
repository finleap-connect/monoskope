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

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	aggregates "github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// UserRoleBindingAggregate is an aggregate for UserRoleBindings.
type UserRoleBindingAggregate struct {
	*DomainAggregateBase
	aggregateManager es.AggregateStore
	userId           uuid.UUID // User to add a role to
	role             es.Role   // Role to add to the user
	scope            es.Scope  // Scope of the role binding
	resource         uuid.UUID // Resource of the role binding
}

// NewUserRoleBindingAggregate creates a new UserRoleBindingAggregate
func NewUserRoleBindingAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &UserRoleBindingAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.UserRoleBinding),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.validate(ctx, cmd); err != nil {
		return nil, err
	}
	return a.execute(ctx, cmd)
}

func (a *UserRoleBindingAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserRoleBindingCommand:
		if a.Exists() {
			return domainErrors.ErrUserRoleBindingAlreadyExists
		}

		var err error
		var userId uuid.UUID
		resource := uuid.Nil
		// Get all aggregates of same type
		if userId, err = uuid.Parse(cmd.GetUserId()); err != nil {
			return domainErrors.ErrInvalidArgument("user id is invalid")
		}
		if err := roles.ValidateRole(cmd.GetRole()); err != nil {
			return err
		}
		if err := scopes.ValidateScope(cmd.GetScope()); err != nil {
			return err
		}
		if cmd.Resource != nil {
			resourceValue := cmd.GetResource().GetValue()
			if resource, err = uuid.Parse(resourceValue); err != nil {
				return domainErrors.ErrInvalidArgument("resource id is invalid")
			}
		}

		userAggregate, err := a.aggregateManager.Get(ctx, aggregates.User, userId)
		if err != nil {
			return err
		}
		if !userAggregate.Exists() || userAggregate.Deleted() {
			return domainErrors.ErrUserNotFound
		}

		roleBindings, err := a.aggregateManager.All(ctx, aggregates.UserRoleBinding)
		if err != nil {
			return err
		}
		if containsRoleBinding(roleBindings, cmd.UserId, cmd.Role, cmd.Scope, resource.String()) {
			return domainErrors.ErrUserRoleBindingAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func (a *UserRoleBindingAggregate) execute(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	switch cmd := cmd.(type) {
	case *commands.CreateUserRoleBindingCommand:
		resource := ""
		resourceValue := cmd.GetResource()
		if resourceValue != nil {
			resource = resourceValue.Value
		}
		eventData := &eventdata.UserRoleAdded{
			UserId:   cmd.GetUserId(),
			Role:     cmd.GetRole(),
			Scope:    cmd.GetScope(),
			Resource: resource,
		}
		_ = a.AppendEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(eventData))
	case *commands.DeleteUserRoleBindingCommand:
		_ = a.AppendEvent(ctx, events.UserRoleBindingDeleted, nil)
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	return a.DefaultReply(), nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserRoleBindingCreated:
		err := a.userRoleBindingCreated(event)
		if err != nil {
			return err
		}
	case events.UserRoleBindingDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}

// userRoleBindingCreated applies the event on the aggregate
func (a *UserRoleBindingAggregate) userRoleBindingCreated(event es.Event) error {
	data := &eventdata.UserRoleAdded{}
	err := event.Data().ToProto(data)
	if err != nil {
		return err
	}

	userId, err := uuid.Parse(data.UserId)
	if err != nil {
		return err
	}

	a.userId = userId
	a.role = es.Role(data.Role)
	a.scope = es.Scope(data.Scope)
	a.resource = uuid.Nil

	if data.Resource != "" {
		id, err := uuid.Parse(data.Resource)
		if err != nil {
			return err
		}
		a.resource = id
	}

	return nil
}

func containsRoleBinding(values []es.Aggregate, userId string, role, scope, resource string) bool {
	resourceId := uuid.Nil
	if resource != "" {
		id, err := uuid.Parse(resource)
		if err != nil {
			return false
		}
		resourceId = id
	}

	for _, value := range values {
		d, ok := value.(*UserRoleBindingAggregate)
		if ok &&
			d.userId.String() == userId &&
			d.role.String() == role &&
			d.scope.String() == scope &&
			d.resource == resourceId {
			return true
		}
	}
	return false
}
