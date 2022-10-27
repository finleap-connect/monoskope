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
	for _, eventType := range events.ClusterEvents {
		_ = event.DefaultEventFormatterRegistry.RegisterEventFormatter(eventType, NewClusterEventFormatter)
	}
}

// clusterEventFormatter EventFormatter implementation for the cluster-aggregate
type clusterEventFormatter struct {
	*event.EventFormatterBase
	esClient esApi.EventStoreClient
}

// NewClusterEventFormatter creates a new event formatter for the cluster-aggregate
func NewClusterEventFormatter(esClient esApi.EventStoreClient) event.EventFormatter {
	return &clusterEventFormatter{
		&event.EventFormatterBase{}, esClient,
	}
}

// GetFormattedDetails formats the cluster-aggregate-events in a human-readable format
func (f *clusterEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.ClusterDeleted:
		return f.getFormattedDetailsClusterDeleted(ctx, event)
	}

	ed, err := es.EventData(event.Data).Unmarshal()
	if err != nil {
		return "", err
	}

	switch ed := ed.(type) {
	case *eventdata.ClusterCreated:
		return f.getFormattedDetailsClusterCreated(event, ed)
	case *eventdata.ClusterCreatedV2:
		return f.getFormattedDetailsClusterCreatedV2(event, ed)
	case *eventdata.ClusterCreatedV3:
		return f.getFormattedDetailsClusterCreatedV3(event, ed)
	case *eventdata.ClusterUpdated:
		return f.getFormattedDetailsClusterUpdated(ctx, event, ed)
	case *eventdata.ClusterUpdatedV2:
		return f.getFormattedDetailsClusterUpdatedV2(ctx, event, ed)
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreated(event *esApi.Event, eventData *eventdata.ClusterCreated) (string, error) {
	return fConsts.ClusterCreatedV2DetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], eventData.Name), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreatedV2(event *esApi.Event, eventData *eventdata.ClusterCreatedV2) (string, error) {
	return fConsts.ClusterCreatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], eventData.Name), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreatedV3(event *esApi.Event, eventData *eventdata.ClusterCreatedV3) (string, error) {
	return fConsts.ClusterCreatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], eventData.Name), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.ClusterUpdated) (string, error) {
	clusterSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewClusterProjector())

	cluster, err := clusterSnapshotter.CreateSnapshot(ctx, &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}

	var details strings.Builder
	details.WriteString(fConsts.ClusterUpdatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]))
	f.AppendUpdate("Display name", eventData.DisplayName, cluster.DisplayName, &details)
	f.AppendUpdate("API server address", eventData.ApiServerAddress, cluster.ApiServerAddress, &details)
	if len(eventData.CaCertificateBundle) != 0 {
		f.AppendUpdate("Certificate", "a new one", "", &details)
	}
	return details.String(), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterUpdatedV2(ctx context.Context, event *esApi.Event, eventData *eventdata.ClusterUpdatedV2) (string, error) {
	snapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewClusterProjector())

	cluster, err := snapshotter.CreateSnapshot(ctx, &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}

	var details strings.Builder
	details.WriteString(fConsts.ClusterUpdatedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]))
	if eventData.Name != nil {
		f.AppendUpdate("Name", eventData.Name.Value, cluster.DisplayName, &details)
	}
	if eventData.ApiServerAddress != nil {
		f.AppendUpdate("API server address", eventData.ApiServerAddress.Value, cluster.ApiServerAddress, &details)
	}
	if len(eventData.CaCertificateBundle) != 0 {
		f.AppendUpdate("Certificate", "a new one", "", &details)
	}
	return details.String(), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	clusterSnapshotter := snapshots.NewSnapshotter(f.esClient, projectors.NewClusterProjector())

	cluster, err := clusterSnapshotter.CreateSnapshot(ctx, &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}

	return fConsts.ClusterDeletedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail], cluster.DisplayName), nil
}
