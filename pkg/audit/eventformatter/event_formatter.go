package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
	"time"
)


type eventFormatter interface {
	getFormattedDetails() string
}

type EventFormatter struct {
	QHDomain domain.QueryHandlerDomain
}

func NewEventFormatter(qhDomain domain.QueryHandlerDomain) *EventFormatter {
	return &EventFormatter{QHDomain: qhDomain}
}

func (f *EventFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	formatter := f.getFormatterBasedOnEventType(ctx, event)

	return &audit.HumanReadableEvent{
		When: event.Timestamp.AsTime().Format(time.RFC822),
		Issuer: event.Metadata["x-auth-email"],
		IssuerId: event.AggregateId,
		EventType: event.Type,
		Details: formatter.getFormattedDetails(),
	}
}

func (f *EventFormatter) getFormatterBasedOnEventType(ctx context.Context, event *esApi.Event) eventFormatter {
	switch es.EventType(event.Type) {
	case events.UserCreated, events.UserDeleted,
	events.UserRoleBindingCreated, events.UserRoleBindingDeleted:
		return newUserEventFormatter(*f, ctx, event)

	case events.ClusterCreated, events.ClusterCreatedV2, events.ClusterUpdated, events.ClusterDeleted,
	events.ClusterBootstrapTokenCreated:
		return newClusterEventFormatter(*f, ctx, event)

	case events.TenantCreated, events.TenantDeleted, events.TenantUpdated,
	events.TenantClusterBindingCreated, events.TenantClusterBindingDeleted:
		return newTenantEventFormatter(*f, ctx, event)

	case events.CertificateRequested, events.CertificateRequestIssued, events.CertificateIssued,
	events.CertificateIssueingFailed:
		return newCertificateEventFormatter(*f, ctx, event)
	}

	return nil
}

func (f *EventFormatter) getUserById(ctx context.Context, id string) (*projections.User, error) {
	aggregateId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := f.QHDomain.UserRepository.ByUserId(ctx, aggregateId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (f *EventFormatter) getUserRoleBindingById(ctx context.Context, id string) (*projections.UserRoleBinding, error) {
	aggregateId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := f.QHDomain.UserRoleBindingRepository.ById(ctx, aggregateId)
	if err != nil {
		return nil, err
	}

	urb, ok := projection.(*projections.UserRoleBinding)
	if !ok {
		return nil, errors.ErrProjectionNotFound
	}
	
	return urb, nil
}

func (f *EventFormatter) getClusterById(ctx context.Context, id string) (*projections.Cluster, error) {
	cluster, err := f.QHDomain.ClusterRepository.ByClusterId(ctx, id)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (f *EventFormatter) getTenantById(ctx context.Context, id string) (*projections.Tenant, error) {
	tenant, err := f.QHDomain.TenantRepository.ByTenantId(ctx, id)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (f *EventFormatter) getTenantClusterBinding(ctx context.Context, id string) (*projections.TenantClusterBinding, error) {
	aggregateId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := f.QHDomain.TenantClusterBindingRepository.ById(ctx, aggregateId)
	if err != nil {
		return nil, err
	}

	tcb, ok := projection.(*projections.TenantClusterBinding)
	if !ok {
		return nil, errors.ErrProjectionNotFound
	}

	return tcb, nil
}

func appendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf("\n- %s to %s", field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(" from %s", old))
		}
	}
}

// TODO: make this public in event data util
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