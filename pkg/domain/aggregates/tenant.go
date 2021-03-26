package aggregates

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// TenantAggregate is an aggregate for Tenants.
type TenantAggregate struct {
	*es.BaseAggregate
	tenantRepo repos.ReadOnlyTenantRepository
	Name       string
	Prefix     string
}

// NewTenantAggregate creates a new TenantAggregate
func NewTenantAggregate(id uuid.UUID, tenantRepo repos.ReadOnlyTenantRepository) *TenantAggregate {
	return &TenantAggregate{
		BaseAggregate: es.NewBaseAggregate(aggregates.Tenant, id),
		tenantRepo:    tenantRepo,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *TenantAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateTenantCommand:
		_, err := a.tenantRepo.ByName(ctx, cmd.GetName())
		if err != nil && errors.Is(err, domainErrors.ErrTenantNotFound) {
			ed := es.ToEventDataFromProto(&eventdata.TenantCreatedEventData{Name: cmd.GetName(), Prefix: cmd.GetPrefix()})
			_ = a.AppendEvent(ctx, events.TenantCreated, ed)
			return nil
		} else {
			return domainErrors.ErrTenantAlreadyExists
		}
	case *commands.UpdateTenantCommand:
		// Check if tenant exists
		if err := a.checkExists(); err != nil {
			return err
		}

		ed := es.ToEventDataFromProto(&eventdata.TenantUpdatedEventData{
			Id:     cmd.GetId(),
			Update: &eventdata.TenantUpdatedEventData_Update{Name: cmd.GetUpdate().GetName()},
		})
		_ = a.AppendEvent(ctx, events.TenantUpdated, ed)
		return nil
	case *commands.DeleteTenantCommand:
		// Check if tenant exists
		if err := a.checkExists(); err != nil {
			return err
		}

		ed := es.ToEventDataFromProto(&eventdata.TenantDeletedEventData{Id: cmd.GetId()})
		_ = a.AppendEvent(ctx, events.TenantDeleted, ed)
		return nil
	}

	return fmt.Errorf("couldn't handle command")
}

func (a *TenantAggregate) checkExists() error {
	if a.Version() < 1 {
		return fmt.Errorf("tenant does not exist")
	}
	return nil
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *TenantAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.TenantCreated:
		data := &eventdata.TenantCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Name = data.Name
		a.Prefix = data.Prefix
	case events.TenantUpdated:
		data := &eventdata.TenantUpdatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return err
		}
		a.Name = data.Update.Name.Value
	case events.TenantDeleted:
		a.SetDeleted()
	default:
		return fmt.Errorf("couldn't handle event")
	}
	return nil
}
