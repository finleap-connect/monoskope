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

	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	aggregates "github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	*DomainAggregateBase
	aggregateManager es.AggregateStore
	Email            string
	Name             string
}

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &UserAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.User),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.Authorize(ctx, cmd, uuid.Nil); err != nil {
		return nil, err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return nil, err
	}
	return a.execute(ctx, cmd)
}

func (a *UserAggregate) execute(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		source, err := sourceFromContext(ctx)
		if err != nil {
			return nil, err
		}
		_ = a.AppendEvent(ctx, events.UserCreated, es.ToEventDataFromProto(&eventdata.UserCreated{
			Email:  cmd.GetEmail(),
			Name:   cmd.GetName(),
			Source: source,
		}))
		reply := &es.CommandReply{
			Id:      a.ID(),
			Version: a.Version(),
		}
		return reply, nil
	case *commands.UpdateUserCommand:
		ed := new(eventdata.UserUpdated)
		name := cmd.GetName()

		if name != nil && a.Name != name.Value {
			ed.Name = name.Value
			_ = a.AppendEvent(ctx, events.UserUpdated, es.ToEventDataFromProto(ed))
		}
		reply := &es.CommandReply{
			Id:      a.ID(),
			Version: a.Version(),
		}
		return reply, nil
	case *commands.DeleteUserCommand:
		_ = a.AppendEvent(ctx, events.UserDeleted, nil)
		reply := &es.CommandReply{
			Id:      a.ID(),
			Version: a.Version(),
		}
		return reply, nil
	}
	return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
}

func (a *UserAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		if a.Exists() {
			return domainErrors.ErrUserAlreadyExists
		}

		// Get all aggregates of same type
		aggregates, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}

		// Check if user already exists
		if containsUser(aggregates, cmd.GetEmail()) {
			return domainErrors.ErrUserAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserCreated:
		err := a.userCreated(event)
		if err != nil {
			return err
		}
	case events.UserUpdated:
		data := new(eventdata.UserUpdated)
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.Name = data.GetName()
	case events.UserDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}

	return nil
}

// userCreated handle the event
func (a *UserAggregate) userCreated(event es.Event) error {
	data := &eventdata.UserCreated{}
	if err := event.Data().ToProto(data); err != nil {
		return err
	}

	a.Email = data.Email
	a.Name = data.Name

	return nil
}

func containsUser(values []es.Aggregate, emailAddress string) bool {
	for _, value := range values {
		d, ok := value.(*UserAggregate)
		if ok {
			if !d.Deleted() && d.Email == emailAddress {
				return true
			}
		}
	}
	return false
}

func sourceFromContext(ctx context.Context) (common.UserSource, error) {
	// Extract domain context
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return common.UserSource_INTERNAL, err
	}

	switch metadataManager.GetComponentName() {
	case users.SCIMServerUser.Name:
		return common.UserSource_SCIM, nil
	default:
		return common.UserSource_INTERNAL, nil
	}
}
