package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"strings"
	"time"
)


type eventFormatter interface {
	getFormattedDetails(context.Context) string
}

type EventFormatter struct {
	esClient esApi.EventStoreClient
}

func NewEventFormatter(esClient esApi.EventStoreClient) *EventFormatter {
	return &EventFormatter{esClient: esClient}
}

func (f *EventFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	formatter := f.getFormatterBasedOnEventType(event)

	return &audit.HumanReadableEvent{
		When: event.Timestamp.AsTime().Format(time.RFC822),
		Issuer: event.Metadata["x-auth-email"],
		IssuerId: event.AggregateId,
		EventType: event.Type,
		Details: formatter.getFormattedDetails(ctx),
	}
}

func (f *EventFormatter) getFormatterBasedOnEventType(event *esApi.Event) eventFormatter {
	switch es.EventType(event.Type) {
	case events.UserCreated, events.UserDeleted,
		 events.UserRoleBindingCreated, events.UserRoleBindingDeleted:
		return newUserEventFormatter(*f, event)

	case events.TenantCreated, events.TenantDeleted, events.TenantUpdated,
		 events.TenantClusterBindingCreated, events.TenantClusterBindingDeleted:
		return newTenantEventFormatter(*f, event)

	case events.ClusterCreated, events.ClusterCreatedV2, events.ClusterUpdated, events.ClusterDeleted,
		 events.ClusterBootstrapTokenCreated:
		return newClusterEventFormatter(*f, event)

	case events.CertificateRequested, events.CertificateRequestIssued, events.CertificateIssued,
		 events.CertificateIssueingFailed:
		return newCertificateEventFormatter(*f, event)
	}

	return nil
}

// TODO: find a better place
// TODO: ticket: domain -> snapshots (same idea as e.g. domain -> projections)
func (f *EventFormatter) getSnapshot(ctx context.Context, projector es.Projector, eventFilter *esApi.EventFilter) (es.Projection, error) {
	projection := projector.NewProjection(uuid.New())
	aggregateEvents, err := f.esClient.Retrieve(ctx, eventFilter)
	if err != nil {
		return nil, err
	}

	for {
		e, err := aggregateEvents.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			return nil, err
		}

		event, err := es.NewEventFromProto(e)
		if err != nil {
			return nil, err
		}

		projection, err = projector.Project(ctx, event, projection)
		if err != nil {
			return nil, err
		}

	}

	return projection, nil
}

func appendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf("\n- “%s“ to “%s“", field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(" from “%s“", old))
		}
	}
}

func toPortoFromEventData(eventData []byte) (proto.Message, bool) {
	porto := &anypb.Any{}
	if err := protojson.Unmarshal(eventData, porto); err != nil {
		return nil, false
	}
	ed, err := porto.UnmarshalNew()
	if err != nil {
		return nil, false
	}
	return ed, true
}