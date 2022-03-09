package formatters

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)


type tenantEventFormatter struct {
	*eventformatter.BaseEventFormatter
}

func NewTenantEventFormatter(esClient esApi.EventStoreClient) *tenantEventFormatter {
	return &tenantEventFormatter{
		BaseEventFormatter: &eventformatter.BaseEventFormatter{EsClient: esClient},
	}
}

func (f *tenantEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.TenantDeleted: return f.getFormattedDetailsTenantDeleted(ctx, event)
	case events.TenantClusterBindingDeleted: return f.getFormattedDetailsTenantClusterBindingDeleted(ctx, event)
	}

	ed, err := f.ToPortoFromEventData(event.Data)
	if err != nil {
		return "", err
	}

	switch ed.(type) {
	case *eventdata.TenantCreated: return f.getFormattedDetailsTenantCreated(event, ed.(*eventdata.TenantCreated))
	case *eventdata.TenantUpdated: return f.getFormattedDetailsTenantUpdated(ctx, event, ed.(*eventdata.TenantUpdated))
	case *eventdata.TenantClusterBindingCreated: return f.getFormattedDetailsTenantClusterBindingCreated(ctx, event, ed.(*eventdata.TenantClusterBindingCreated))
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *tenantEventFormatter) getFormattedDetailsTenantCreated(event *esApi.Event, eventData *eventdata.TenantCreated) (string, error) {
	return fmt.Sprintf("“%s“ created tenant “%s“ with prefix “%s“", event.Metadata["x-auth-email"], eventData.Name, eventData.Prefix), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingCreated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantClusterBindingCreated) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.TenantId}
	tenantSnapshot, err := f.GetSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.ClusterId}
	clusterSnapshot, err := f.GetSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ bounded tenant “%s“ to cluster “%s”",
		event.Metadata["x-auth-email"], tenant.Name, cluster.DisplayName), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantUpdated) (string, error) {
	tenantSnapshot, err := f.GetSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId: &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	oldTenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the Tenant", event.Metadata["x-auth-email"]))
	f.AppendUpdate("Name", eventData.Name.Value, oldTenant.Name, &details)
	return details.String(), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	tenantSnapshot, err := f.GetSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ deleted tenant “%s“", event.Metadata["x-auth-email"], tenant.Name), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: event.AggregateId}
	tcbSnapshot, err := f.GetSnapshot(ctx, projectors.NewTenantClusterBindingProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tcb, ok := tcbSnapshot.(*projections.TenantClusterBinding)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.TenantId}
	tenantSnapshot, err := f.GetSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.ClusterId}
	clusterSnapshot, err := f.GetSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ deleted the bound between cluster “%s“ and tenant “%s“",
		event.Metadata["x-auth-email"], cluster.DisplayName, tenant.Name), nil
}