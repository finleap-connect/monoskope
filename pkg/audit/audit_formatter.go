package audit

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"time"
)

type AuditFormatter interface {
	NewHumanReadableEvent(context.Context, *esApi.Event) *audit.HumanReadableEvent
}

type auditFormatter struct {
	log logger.Logger
	eventFormatterRegistry eventformatter.EventFormatterRegistry
}

func NewAuditFormatter(esClient esApi.EventStoreClient) *auditFormatter {
	return &auditFormatter{
		log: logger.WithName("audit-formatter"),
		eventFormatterRegistry: initEventFormatterRegistry(esClient),
	}
}

func (f *auditFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	return &audit.HumanReadableEvent{
		When: event.Timestamp.AsTime().Format(time.RFC822),
		Issuer: event.Metadata["x-auth-email"],
		IssuerId: event.AggregateId,
		EventType: event.Type,
		Details: f.getFormattedDetails(ctx, event),
	}
}

func (f *auditFormatter) getFormattedDetails(ctx context.Context, event *esApi.Event) string {
	eventFormatter, err := f.eventFormatterRegistry.GetEventFormatter(es.EventType(event.Type))
	if err != nil {
		return ""
	}

	details, err := eventFormatter.GetFormattedDetails(ctx, event)
	if err != nil {
		f.log.Error(err, "failed to format event details",
			"eventAggregate", event.GetAggregateId(),
			"eventTimestamp", event.GetTimestamp().AsTime().Format(time.RFC3339Nano))
	}

	return details
}


// TODO: default registry and init on start?
func initEventFormatterRegistry(esClient esApi.EventStoreClient) eventformatter.EventFormatterRegistry {
	eventFormatterRegistry := eventformatter.NewEventFormatterRegistry()

	// TODO: group/enum

	// User
	userEvents := [...]es.EventType{events.UserCreated, events.UserDeleted,
		events.UserRoleBindingCreated, events.UserRoleBindingDeleted}
	userEventsFormatter := formatters.NewUserEventFormatter(esClient)
	for _, eventType := range userEvents {
		_ = eventFormatterRegistry.RegisterEventFormatter(eventType, userEventsFormatter)
	}

	// Tenant
	tenantEvents := [...]es.EventType{events.TenantCreated, events.TenantDeleted, events.TenantUpdated,
		events.TenantClusterBindingCreated, events.TenantClusterBindingDeleted}
	tenantEventsFormatter := formatters.NewTenantEventFormatter(esClient)
	for _, eventType := range tenantEvents {
		_ = eventFormatterRegistry.RegisterEventFormatter(eventType, tenantEventsFormatter)
	}

	// Cluster
	clusterEvents := [...]es.EventType{events.ClusterCreated, events.ClusterCreatedV2, events.ClusterUpdated, events.ClusterDeleted,
		events.ClusterBootstrapTokenCreated}
	clusterEventsFormatter := formatters.NewClusterEventFormatter(esClient)
	for _, eventType := range clusterEvents {
		_ = eventFormatterRegistry.RegisterEventFormatter(eventType, clusterEventsFormatter)
	}

	// Certificate
	certificateEvents := [...]es.EventType{events.CertificateRequested, events.CertificateRequestIssued, events.CertificateIssued,
		events.CertificateIssueingFailed}
	certificateEventsFormatter := formatters.NewCertificateEventFormatter(esClient)
	for _, eventType := range certificateEvents {
		_ = eventFormatterRegistry.RegisterEventFormatter(eventType, certificateEventsFormatter)
	}

	return eventFormatterRegistry
}