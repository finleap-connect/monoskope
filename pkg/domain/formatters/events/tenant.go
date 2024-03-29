// Copyright 2022 Monoskope Authors
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

package events

import (
	"context"
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	fConsts "github.com/finleap-connect/monoskope/pkg/domain/constants/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"github.com/finleap-connect/monoskope/pkg/domain/snapshots"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func init() {
	for _, eventType := range events.TenantEvents {
		_ = event.DefaultEventFormatterRegistry.RegisterEventFormatter(eventType, NewTenantEventFormatter)
	}
}

// tenantEventFormatter EventFormatter implementation for the tenant-aggregate
type tenantEventFormatter struct {
	*event.EventFormatterBase
	esClient esApi.EventStoreClient
}

// NewTenantEventFormatter creates a new event formatter for the tenant-aggregate
func NewTenantEventFormatter(esClient esApi.EventStoreClient) event.EventFormatter {
	return &tenantEventFormatter{
		&event.EventFormatterBase{}, esClient,
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
	return fConsts.TenantCreatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], eventData.Name, eventData.Prefix), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantUpdated) (string, error) {
	tenantSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewTenantProjector())

	tenant, err := tenantSnapshotter.CreateSnapshot(ctx, &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId},
	})
	if err != nil {
		return "", err
	}

	var details strings.Builder
	details.WriteString(fConsts.TenantUpdatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]))
	f.AppendUpdate("Name", eventData.Name.Value, tenant.Name, &details)
	return details.String(), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingCreated(ctx context.Context, event *esApi.Event, eventData *eventdata.TenantClusterBindingCreated) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.TenantId}

	tenantSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewTenantProjector())
	tenant, err := tenantSnapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		return "", err
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: eventData.ClusterId}

	clusterSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewClusterProjector())
	cluster, err := clusterSnapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		return "", err
	}

	return fConsts.TenantClusterBindingCreatedDetailsFormat.Sprint(
		event.Metadata[auth.HeaderAuthEmail], tenant.Name, cluster.Name), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	tenantSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewTenantProjector())
	tenant, err := tenantSnapshotter.CreateSnapshot(ctx, &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId},
	})
	if err != nil {
		return "", err
	}

	return fConsts.TenantDeletedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], tenant.Name), nil
}

func (f *tenantEventFormatter) getFormattedDetailsTenantClusterBindingDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: event.AggregateId}

	tenantClusterBindingSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewTenantClusterBindingProjector())
	tcb, err := tenantClusterBindingSnapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		return "", err
	}

	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.TenantId}
	tenantSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewTenantProjector())
	tenant, err := tenantSnapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		return "", err
	}

	eventFilter.AggregateId = &wrapperspb.StringValue{Value: tcb.ClusterId}
	clusterSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewClusterProjector())
	cluster, err := clusterSnapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		return "", err
	}

	return fConsts.TenantClusterBindingDeletedDetailsFormat.Sprint(
		event.Metadata[auth.HeaderAuthEmail], cluster.Name, tenant.Name), nil
}
