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
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// TenantAggregate is an aggregate for Tenants.
type TenantAggregate struct {
	DomainAggregateBase
	aggregateManager es.AggregateManager
	Name             string
	Prefix           string
}

// NewTenantAggregate creates a new TenantAggregate
func NewTenantAggregate(id uuid.UUID, aggregateManager es.AggregateManager) es.Aggregate {
	return &TenantAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Tenant, id),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *TenantAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return err
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
			if d.Name == name {
				return true
			}
		}
	}
	return false
}

// execute executes the command after it has successfully been validated
func (a *TenantAggregate) execute(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateTenantCommand:
		ed := es.ToEventDataFromProto(&eventdata.TenantCreated{Name: cmd.GetName(), Prefix: cmd.GetPrefix()})
		_ = a.AppendEvent(ctx, events.TenantCreated, ed)
		return nil
	case *commands.UpdateTenantCommand:
		ed := es.ToEventDataFromProto(&eventdata.TenantUpdated{
			Name: cmd.GetName(),
		})
		_ = a.AppendEvent(ctx, events.TenantUpdated, ed)
		return nil
	case *commands.DeleteTenantCommand:
		_ = a.AppendEvent(ctx, events.TenantDeleted, nil)
		return nil
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *TenantAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.TenantCreated:
		data := &eventdata.TenantCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Name = data.Name
		a.Prefix = data.Prefix
	case events.TenantUpdated:
		data := &eventdata.TenantUpdated{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Name = data.GetName().GetValue()
	case events.TenantDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
