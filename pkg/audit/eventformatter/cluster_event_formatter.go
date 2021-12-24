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


type clusterEventFormatter struct {
	EventFormatter
	ctx   context.Context
	event *esApi.Event
}

func newClusterEventFormatter(eventFormatter EventFormatter, ctx context.Context, event *esApi.Event) *clusterEventFormatter {
	return &clusterEventFormatter{EventFormatter: eventFormatter, ctx: ctx, event: event}
}

func (f *clusterEventFormatter) getFormattedDetails() string {
	switch es.EventType(f.event.Type) {
	case events.ClusterDeleted: return f.getFormattedDetailsClusterDeleted()
	}

	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}

	switch ed.(type) {
	case *eventdata.ClusterCreated: return f.getFormattedDetailsClusterCreated(ed.(*eventdata.ClusterCreated))
	case *eventdata.ClusterCreatedV2: return f.getFormattedDetailsClusterCreatedV2(ed.(*eventdata.ClusterCreatedV2))
	case *eventdata.ClusterBootstrapTokenCreated: return f.getFormattedDetailsClusterBootstrapTokenCreated(ed.(*eventdata.ClusterBootstrapTokenCreated))
	case *eventdata.ClusterUpdated: return f.getFormattedDetailsClusterUpdated(ed.(*eventdata.ClusterUpdated))
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

func (f *clusterEventFormatter) getFormattedDetailsClusterUpdated(eventData *eventdata.ClusterUpdated) string {
	// TODO: how to get a projection of a specific version
	oldCluster, err := f.QHDomain.ClusterRepository.ByClusterId(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the cluster", f.event.Metadata["x-auth-email"]))
	appendUpdate("Display name", eventData.DisplayName, oldCluster.DisplayName, &details)
	appendUpdate("API server address", eventData.ApiServerAddress, oldCluster.ApiServerAddress, &details)
	if len(eventData.CaCertificateBundle) != 0 {details.WriteString(fmt.Sprintf("\n- Certifcate to a new one"))}
	return details.String()
}

func (f *clusterEventFormatter) getFormattedDetailsClusterDeleted() string {
	cluster, err := f.QHDomain.ClusterRepository.ByClusterId(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("“%s“ deleted cluster “%s“", f.event.Metadata["x-auth-email"], cluster.DisplayName)
}