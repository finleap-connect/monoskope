package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// UserAggregate is an aggregate for Users.
type UserAggregate struct {
	DomainAggregateBase
	aggregateManager es.AggregateStore
	Email            string
	Name             string
}

// NewUserAggregate creates a new UserAggregate
func NewUserAggregate(id uuid.UUID, aggregateManager es.AggregateStore) es.Aggregate {
	return &UserAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.User, id),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return err
	}
	return a.execute(ctx, cmd)
}

func (a *UserAggregate) execute(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateUserCommand:
		_ = a.AppendEvent(ctx, events.UserCreated, es.ToEventDataFromProto(&eventdata.UserCreated{
			Email: cmd.GetEmail(),
			Name:  cmd.GetName(),
		}))
		return nil
	}
	return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
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

func (a *UserAggregate) authorize(ctx context.Context, cmd es.Command) error {
	metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	if metadataMgr.IsAuthorizationBypassed() {
		return nil
	}

	userRoleBindings := metadataMgr.GetRoleBindings()
	for _, policy := range cmd.Policies(ctx) {
		for _, roleBinding := range userRoleBindings {
			if policy.AcceptsRole(es.Role(roleBinding.Role)) &&
				policy.AcceptsScope(es.Scope(roleBinding.Scope)) {
				return nil
			}
		}
	}
	return domainErrors.ErrUnauthorized
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *UserAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.UserCreated:
		err := a.userCreated(event)
		if err != nil {
			return err
		}
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
			if d.Email == emailAddress {
				return true
			}
		}
	}
	return false
}
