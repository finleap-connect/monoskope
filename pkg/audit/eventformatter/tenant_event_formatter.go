package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)


type tenantEventFormatter struct {
	EventFormatter
	event *esApi.Event
}

func newTenantEventFormatter(eventFormatter EventFormatter, event *esApi.Event) *tenantEventFormatter {
	return &tenantEventFormatter{EventFormatter: eventFormatter, event: event}
}

func (f *tenantEventFormatter) getFormattedDetails(ctx context.Context) string {
	switch es.EventType(f.event.Type) {
	case events.TenantDeleted: return f.getFormattedDetailsTenantDeleted(ctx)
	case events.TenantClusterBindingDeleted: return f.getFormattedDetailsTenantClusterBindingDeleted(ctx)
	}

	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}
	
	switch ed.(type) {
	case *eventdata.TenantCreated: return f.getFormattedDetailsTenantCreated(ed.(*eventdata.TenantCreated))
	case *eventdata.TenantUpdated: return f.getFormattedDetailsTenantUpdated(ctx, ed.(*eventdata.TenantUpdated))
	case *eventdata.TenantClusterBindingCreated: return f.getFormattedDetailsTenantClusterBindingCreated(ctx, ed.(*eventdata.TenantClusterBindingCreated))
	}

	return ""
}

func (f *tenantEventFormatter) getFormattedDetailsTenantCreated(eventData *eventdata.TenantCreated) string {
	return fmt.Sprintf("“%s“ created tenant “%s“ with prefix “%s“", f.event.Metadata["x-auth-email"], eventData.Name, eventData.Prefix)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingCreated(ctx context.Context, eventData *eventdata.TenantClusterBindingCreated) string {
	eventFilter := &esApi.EventFilter{MaxTimestamp: f.event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.TenantId}
	tenantSnapshot, err := f.getSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if err != nil || !ok {
		return ""
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.ClusterId}
	clusterSnapshot, err := f.getSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ bounded tenant “%s“ to cluster “%s”",
		f.event.Metadata["x-auth-email"], tenant.Name, cluster.DisplayName)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantUpdated(ctx context.Context, eventData *eventdata.TenantUpdated) string {
	tenantSnapshot, err := f.getSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(f.event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId: &wrapperspb.StringValue{Value: f.event.AggregateId}},
	)
	oldTenant, ok := tenantSnapshot.(*projections.Tenant)
	if err != nil || !ok {
		return ""
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the Tenant", f.event.Metadata["x-auth-email"]))
	appendUpdate("Name", eventData.Name.Value, oldTenant.Name, &details)
	return details.String()
}

func (f *tenantEventFormatter) getFormattedDetailsTenantDeleted(ctx context.Context) string {
	tenantSnapshot, err := f.getSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: f.event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: f.event.AggregateId}},
	)
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ deleted tenant “%s“", f.event.Metadata["x-auth-email"], tenant.Name)
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingDeleted(ctx context.Context) string {
	eventFilter := &esApi.EventFilter{MaxTimestamp: f.event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: f.event.AggregateId}
	tcbSnapshot, err := f.getSnapshot(ctx, projectors.NewTenantClusterBindingProjector(), eventFilter)
	tcb, ok := tcbSnapshot.(*projections.TenantClusterBinding)
	if err != nil || !ok {
		return ""
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.TenantId}
	tenantSnapshot, err := f.getSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if err != nil || !ok {
		return ""
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.ClusterId}
	clusterSnapshot, err := f.getSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ deleted the bound between cluster “%s“ and tenant “%s“",
		f.event.Metadata["x-auth-email"], cluster.DisplayName, tenant.Name)
}