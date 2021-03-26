package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// UserRoleBindingAggregate is an aggregate for UserRoleBindings.
type UserRoleBindingAggregate struct {
	DomainAggregateBase
	aggregateManager es.AggregateManager
	userId           uuid.UUID // User to add a role to
	role             es.Role   // Role to add to the user
	scope            es.Scope  // Scope of the role binding
	resource         uuid.UUID // Resource of the role binding
}

// NewUserRoleBindingAggregate creates a new UserRoleBindingAggregate
func NewUserRoleBindingAggregate(id uuid.UUID, aggregateManager es.AggregateManager) es.Aggregate {
	return &UserRoleBindingAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.UserRoleBinding, id),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserRoleBindingAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return err
	}
	return a.execute(ctx, cmd)
}

func (a *UserRoleBindingAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserRoleBindingCommand:
		if a.Exists() {
			return domainErrors.ErrTenantAlreadyExists
		}

		// Get all aggregates of same type
		userAggregate, err := a.aggregateManager.Get(ctx, aggregates.User, uuid.MustParse(cmd.GetUserId()))
		if err != nil {
			return err
		}

		if userAggregate == nil {
			return domainErrors.ErrUserNotFound
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func (a *UserRoleBindingAggregate) execute(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserRoleBindingCommand:
		eventData := &ed.UserRoleAddedEventData{
			UserId:   cmd.GetUserId(),
			Role:     cmd.GetRole(),
			Scope:    cmd.GetScope(),
			Resource: cmd.GetResource(),
		}
		_ = a.AppendEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(eventData))
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	return nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserRoleBindingCreated:
		err := a.userRoleBindingCreated(event)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}

// userRoleBindingCreated applies the event on the aggregate
func (a *UserRoleBindingAggregate) userRoleBindingCreated(event es.Event) error {
	data := &ed.UserRoleAddedEventData{}
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

	id, err := uuid.Parse(data.Resource)
	if err != nil {
		return err
	}
	a.resource = id

	return nil
}
