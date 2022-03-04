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


type clusterEventFormatter struct {
	EventFormatter
	event *esApi.Event
}

func newClusterEventFormatter(eventFormatter EventFormatter, event *esApi.Event) *clusterEventFormatter {
	return &clusterEventFormatter{EventFormatter: eventFormatter, event: event}
}

func (f *clusterEventFormatter) getFormattedDetails(ctx context.Context) string {
	switch es.EventType(f.event.Type) {
	case events.ClusterDeleted: return f.getFormattedDetailsClusterDeleted(ctx)
	}

	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}

	switch ed.(type) {
	case *eventdata.ClusterCreated: return f.getFormattedDetailsClusterCreated(ed.(*eventdata.ClusterCreated))
	case *eventdata.ClusterCreatedV2: return f.getFormattedDetailsClusterCreatedV2(ed.(*eventdata.ClusterCreatedV2))
	case *eventdata.ClusterBootstrapTokenCreated: return f.getFormattedDetailsClusterBootstrapTokenCreated(ed.(*eventdata.ClusterBootstrapTokenCreated))
	case *eventdata.ClusterUpdated: return f.getFormattedDetailsClusterUpdated(ctx, ed.(*eventdata.ClusterUpdated))
	}

	return ""
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreated(eventData *eventdata.ClusterCreated) string {
	return fmt.Sprintf("“%s“ created cluster “%s“", f.event.Metadata["x-auth-email"], eventData.Name)
}

func (f *clusterEventFormatter) getFormattedDetailsClusterCreatedV2(eventData *eventdata.ClusterCreatedV2) string {
	return fmt.Sprintf("“%s“ created cluster “%s“", f.event.Metadata["x-auth-email"], eventData.Name)
}

func (f *clusterEventFormatter) getFormattedDetailsClusterBootstrapTokenCreated(_ *eventdata.ClusterBootstrapTokenCreated) string {
	return fmt.Sprintf("“%s“ created a cluster bootstrap token", f.event.Metadata["x-auth-email"])
}

func (f *clusterEventFormatter) getFormattedDetailsClusterUpdated(ctx context.Context, eventData *eventdata.ClusterUpdated) string {
	clusterSnapshot, err := f.getSnapshot(ctx, projectors.NewClusterProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(f.event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId: &wrapperspb.StringValue{Value: f.event.AggregateId}},
	)
	oldCluster, ok := clusterSnapshot.(*projections.Cluster)
	if err != nil || !ok {
		return ""
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the cluster", f.event.Metadata["x-auth-email"]))
	appendUpdate("Display name", eventData.DisplayName, oldCluster.DisplayName, &details)
	appendUpdate("API server address", eventData.ApiServerAddress, oldCluster.ApiServerAddress, &details)
	if len(eventData.CaCertificateBundle) != 0 {details.WriteString(fmt.Sprintf("\n- Certifcate to a new one"))}
	return details.String()
}

func (f *clusterEventFormatter) getFormattedDetailsClusterDeleted(ctx context.Context) string {
	clusterSnapshot, err := f.getSnapshot(ctx, projectors.NewClusterProjector(), &esApi.EventFilter{
		MaxTimestamp: f.event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: f.event.AggregateId}},
	)
	cluster, ok := clusterSnapshot.(*projections.Cluster)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ deleted cluster “%s“", f.event.Metadata["x-auth-email"], cluster.DisplayName)
}