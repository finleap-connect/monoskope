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


type clusterEventFormatter struct {
	*eventformatter.BaseEventFormatter
}

func NewClusterEventFormatter(esClient esApi.EventStoreClient) *clusterEventFormatter {
	return &clusterEventFormatter{
		BaseEventFormatter: &eventformatter.BaseEventFormatter{EsClient: esClient},
	}
}

func (f *clusterEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.ClusterDeleted: return f.getFormattedDetailsClusterDeleted(ctx, event)
	}

	ed, err := es.EventData(event.Data).Unmarshal()
	if err != nil {
		return "", err
	}

	switch ed.(type) {
	case *eventdata.ClusterCreated: return f.getFormattedDetailsClusterCreated(event, ed.(*eventdata.ClusterCreated))
	case *eventdata.ClusterCreatedV2: return f.getFormattedDetailsClusterCreatedV2(event, ed.(*eventdata.ClusterCreatedV2))
	case *eventdata.ClusterBootstrapTokenCreated: return f.getFormattedDetailsClusterBootstrapTokenCreated(event)
	case *eventdata.ClusterUpdated: return f.getFormattedDetailsClusterUpdated(ctx, event, ed.(*eventdata.ClusterUpdated))
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreated(event *esApi.Event, eventData *eventdata.ClusterCreated) (string, error) {
	return fmt.Sprintf("“%s“ created cluster “%s“", event.Metadata["x-auth-email"], eventData.Name), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreatedV2(event *esApi.Event, eventData *eventdata.ClusterCreatedV2) (string, error) {
	return fmt.Sprintf("“%s“ created cluster “%s“", event.Metadata["x-auth-email"], eventData.Name), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterBootstrapTokenCreated(event *esApi.Event) (string, error) {
	return fmt.Sprintf("“%s“ created a cluster bootstrap token", event.Metadata["x-auth-email"]), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.ClusterUpdated) (string, error) {
	clusterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewClusterProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId: &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	oldCluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}


	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the cluster", event.Metadata["x-auth-email"]))
	f.AppendUpdate("Display name", eventData.DisplayName, oldCluster.DisplayName, &details)
	f.AppendUpdate("API server address", eventData.ApiServerAddress, oldCluster.ApiServerAddress, &details)
	if len(eventData.CaCertificateBundle) != 0 {details.WriteString(fmt.Sprintf("\n- Certifcate to a new one"))}
	return details.String(), nil
}

func (f *clusterEventFormatter) getFormattedDetailsClusterDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	clusterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewClusterProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ deleted cluster “%s“", event.Metadata["x-auth-email"], cluster.DisplayName), nil
}