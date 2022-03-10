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

const (
	TenantCreatedDetails               = "“%s“ created tenant “%s“ with prefix “%s“"
	TenantUpdatedDetails               = "“%s“ updated the Tenant"
	TenantClusterBindingCreatedDetails = "“%s“ bounded tenant “%s“ to cluster “%s”"
	TenantDeletedDetails               = "“%s“ deleted tenant “%s“"
	TenantClusterBindingDeletedDetails = "“%s“ deleted the bound between cluster “%s“ and tenant “%s“"
)

// tenantEventFormatter EventFormatter implementation for the tenant-aggregate
type tenantEventFormatter struct {
	*eventformatter.BaseEventFormatter
}

// NewTenantEventFormatter creates a new event formatter for the tenant-aggregate
func NewTenantEventFormatter(esClient esApi.EventStoreClient) *tenantEventFormatter {
	return &tenantEventFormatter{
		BaseEventFormatter: &eventformatter.BaseEventFormatter{EsClient: esClient},
	}
}

// GetFormattedDetails formats the tenant-aggregate-events in a human-readable format
func (f *tenantEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.TenantDeleted:
		return f.getFormattedDetailsTenantDeleted(ctx, event)
	case events.TenantClusterBindingDeleted:
		return f.getFormattedDetailsTenantClusterBindingDeleted(ctx, event)
	}

	ed, err := es.EventData(event.Data).Unmarshal()
	if err != nil {
		return "", err
	}

	switch ed := ed.(type) {
	case *eventdata.TenantCreated:
		return f.getFormattedDetailsTenantCreated(event, ed)
	case *eventdata.TenantUpdated:
		return f.getFormattedDetailsTenantUpdated(ctx, event, ed)
	case *eventdata.TenantClusterBindingCreated:
		return f.getFormattedDetailsTenantClusterBindingCreated(ctx, event, ed)
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *tenantEventFormatter) getFormattedDetailsTenantCreated(event *esApi.Event, eventData *eventdata.TenantCreated) (string, error) {
	return fmt.Sprintf(TenantCreatedDetails, event.Metadata["x-auth-email"], eventData.Name, eventData.Prefix), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantUpdated) (string, error) {
	tenantSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	oldTenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf(TenantUpdatedDetails, event.Metadata["x-auth-email"]))
	f.AppendUpdate("Name", eventData.Name.Value, oldTenant.Name, &details)
	return details.String(), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingCreated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantClusterBindingCreated) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.TenantId}
	tenantSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.ClusterId}
	clusterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf(TenantClusterBindingCreatedDetails, event.Metadata["x-auth-email"], tenant.Name, cluster.DisplayName), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	tenantSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf(TenantDeletedDetails, event.Metadata["x-auth-email"], tenant.Name), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: event.AggregateId}
	tcbSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantClusterBindingProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tcb, ok := tcbSnapshot.(*projections.TenantClusterBinding)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.TenantId}
	tenantSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	tenant, ok := tenantSnapshot.(*projections.Tenant)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.ClusterId}
	clusterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf(TenantClusterBindingDeletedDetails, event.Metadata["x-auth-email"], cluster.DisplayName, tenant.Name), nil
}
