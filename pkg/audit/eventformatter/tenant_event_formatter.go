package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"strings"
)


type tenantEventFormatter struct {
	EventFormatter
	ctx   context.Context
	event *esApi.Event
}

func newTenantEventFormatter(eventFormatter EventFormatter, ctx context.Context, event *esApi.Event) *tenantEventFormatter {
	return &tenantEventFormatter{EventFormatter: eventFormatter, ctx: ctx, event: event}
}

func (f *tenantEventFormatter) getFormattedDetails() string {
	switch es.EventType(f.event.Type) {
	case events.TenantDeleted: return f.getFormattedDetailsTenantDeleted()
	case events.TenantClusterBindingDeleted: return f.getFormattedDetailsTenantClusterBindingDeleted()
	}
	
	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}
	switch ed.(type) {
	case *eventdata.TenantCreated: return f.getFormattedDetailsTenantCreated(ed.(*eventdata.TenantCreated))
	case *eventdata.TenantUpdated: return f.getFormattedDetailsTenantUpdated(ed.(*eventdata.TenantUpdated))
	case *eventdata.TenantClusterBindingCreated: return f.getFormattedDetailsTenantClusterBindingCreated(ed.(*eventdata.TenantClusterBindingCreated))
	}

	return ""
}

func (f *tenantEventFormatter) getFormattedDetailsTenantCreated(eventData *eventdata.TenantCreated) string {
	return fmt.Sprintf("%s created Tenant %s with prefix %s", f.event.Metadata["x-auth-email"], eventData.Name, eventData.Prefix)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingCreated(eventData *eventdata.TenantClusterBindingCreated) string {
	tenant, err := f.getTenantById(f.ctx, eventData.TenantId)
	if err != nil {
		return ""
	}
	cluster, err := f.getClusterById(f.ctx, eventData.ClusterId)
	if err != nil {
		return ""
	}
	
	return fmt.Sprintf("%s binded tanent “%s” to cluster “%s”",
		f.event.Metadata["x-auth-email"], tenant.Name, cluster.DisplayName)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantUpdated(eventData *eventdata.TenantUpdated) string {
	// TODO: how to get a projection of a specific version
	oldTenant, err := f.getTenantById(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}
	
	var details strings.Builder
	details.WriteString(fmt.Sprintf("%s updated the Tenant", f.event.Metadata["x-auth-email"]))
	appendUpdate("Name", eventData.Name.Value, oldTenant.Name, &details)
	return details.String()
}

func (f *tenantEventFormatter) getFormattedDetailsTenantDeleted() string {
	tenant, err := f.getTenantById(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}
	
	return fmt.Sprintf("%s deleted tenant %s", f.event.Metadata["x-auth-email"], tenant.Name)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingDeleted() string {
	tcb, err := f.getTenantClusterBinding(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}
	tenant, err := f.getTenantById(f.ctx, tcb.TenantId)
	if err != nil {
		return ""
	}
	cluster, err := f.getClusterById(f.ctx, tcb.ClusterId)
	if err != nil {
		return ""
	}
	
	return fmt.Sprintf("%s deleted the bound between cluster %s and tenant %s",
		f.event.Metadata["x-auth-email"], cluster.DisplayName, tenant.Name)
}