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

package audit

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	ef "github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	_ "github.com/finleap-connect/monoskope/pkg/domain/formatters/events"
	"github.com/finleap-connect/monoskope/pkg/domain/formatters/overviews"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"time"
)

// AuditFormatter is the interface definition for the formatter used by the auditLogServer
type AuditFormatter interface {
	NewHumanReadableEvent(context.Context, *esApi.Event) *audit.HumanReadableEvent
}

// auditFormatter is the implementation of AuditFormatter used by auditLogServer
type auditFormatter struct {
	log        logger.Logger
	esClient   esApi.EventStoreClient
	efRegistry ef.EventFormatterRegistry
	userRepo       repositories.ReadOnlyUserRepository
	tenantRepo     repositories.ReadOnlyTenantRepository
	clusterRepo    repositories.ReadOnlyClusterRepository
}

// NewAuditFormatter creates an auditFormatter
func NewAuditFormatter(esClient esApi.EventStoreClient, efRegistry ef.EventFormatterRegistry) *auditFormatter {
	return &auditFormatter{
		log:        logger.WithName("audit-formatter"),
		esClient:   esClient,
		efRegistry: efRegistry,
	}
}

// NewHumanReadableEvent creates a HumanReadableEvent of a given Event
func (f *auditFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	return &audit.HumanReadableEvent{
		When:      event.Timestamp.AsTime().Format(time.RFC822),
		Issuer:    event.Metadata["x-auth-email"],
		IssuerId:  event.AggregateId,
		EventType: event.Type,
		Details:   f.getEventFormattedDetails(ctx, event),
	}
}

// NewUserOverview TODO
func (f *auditFormatter) NewUserOverview(ctx context.Context, user *projections.User) *audit.UserOverview {
	of := overviews.NewUserOverviewFormatter(f.userRepo, f.tenantRepo, f.clusterRepo)
	uo := &audit.UserOverview{
		Name: user.Name,
		Email: user.Email,
	}
	_ = of.AddRolesDetails(ctx, user, uo)
	uo.Details, _ = of.GetFormattedDetails(ctx, user)
	return uo
}

func (f *auditFormatter) getEventFormattedDetails(ctx context.Context, event *esApi.Event) string {
	eventFormatter, err := f.efRegistry.CreateEventFormatter(f.esClient, es.EventType(event.Type))
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
