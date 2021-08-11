package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
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
		var resource uuid.UUID
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
		if resource, err = uuid.Parse(cmd.GetResource()); err != nil && cmd.GetResource() != "" {
			return domainErrors.ErrInvalidArgument("resource id is invalid")
		}
		if err := a.Authorize(ctx, cmd, resource); err != nil {
			return err
		}

		userAggregate, err := a.aggregateManager.Get(ctx, aggregates.User, userId)
		if err != nil {
			return err
		}
		if !userAggregate.Exists() {
			return domainErrors.ErrUserNotFound
		}

		roleBindings, err := a.aggregateManager.All(ctx, aggregates.UserRoleBinding)
		if err != nil {
			return err
		}
		if containsRoleBinding(roleBindings, cmd.UserId, cmd.Role, cmd.Scope, cmd.Resource) {
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
		eventData := &eventdata.UserRoleAdded{
			UserId:   cmd.GetUserId(),
			Role:     cmd.GetRole(),
			Scope:    cmd.GetScope(),
			Resource: cmd.GetResource(),
		}
		_ = a.AppendEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(eventData))
	case *commands.DeleteUserRoleBindingCommand:
		_ = a.AppendEvent(ctx, events.UserRoleBindingDeleted, nil)
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	reply := &es.CommandReply{
		Id:      a.ID(),
		Version: a.Version(),
	}
	return reply, nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserRoleBindingAggregate) ApplyEvent(event es.Event) error {
	_ = a.BaseAggregate.ApplyEvent(event)

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
