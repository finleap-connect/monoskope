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

package audit

import (
	"context"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	_ "github.com/finleap-connect/monoskope/pkg/domain/formatters/events"
	"github.com/finleap-connect/monoskope/pkg/domain/formatters/overviews"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

// AuditFormatter is the interface definition for the formatter used by the auditLogServer
type AuditFormatter interface {
	NewHumanReadableEvent(context.Context, *esApi.Event) *audit.HumanReadableEvent
	NewUserOverview(context.Context, uuid.UUID, time.Time) *audit.UserOverview
}

// auditFormatter is the implementation of AuditFormatter used by auditLogServer
type auditFormatter struct {
	*formatters.FormatterBase
	log        logger.Logger
	efRegistry event.EventFormatterRegistry
}

// NewAuditFormatter creates an auditFormatter
func NewAuditFormatter(esClient esApi.EventStoreClient, efRegistry event.EventFormatterRegistry) *auditFormatter {
	return &auditFormatter{
		FormatterBase: &formatters.FormatterBase{EsClient: esClient},
		log:           logger.WithName("audit-formatter"),
		efRegistry:    efRegistry,
	}
}

// NewHumanReadableEvent creates a HumanReadableEvent of a given event
func (f *auditFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	humanReadableEvent := &audit.HumanReadableEvent{
		When:      event.Timestamp.AsTime().Format(time.RFC822),
		Issuer:    event.Metadata[auth.HeaderAuthEmail],
		IssuerId:  event.AggregateId,
		EventType: event.Type,
	}

	eventFormatter, err := f.efRegistry.CreateEventFormatter(f.EsClient, es.EventType(event.Type))
	if err != nil {
		return humanReadableEvent
	}
	humanReadableEvent.Details, err = eventFormatter.GetFormattedDetails(ctx, event)
	if err != nil {
		f.log.Error(err, "failed to format event details",
			"eventAggregate", event.GetAggregateId(),
			"eventTimestamp", event.GetTimestamp().AsTime().Format(time.RFC3339Nano))
	}

	return humanReadableEvent
}

// NewUserOverview creates a UserOverview of the given user by its id according to the given timestamp
func (f *auditFormatter) NewUserOverview(ctx context.Context, userId uuid.UUID, timestamp time.Time) *audit.UserOverview {
	userOverview := &audit.UserOverview{}

	userSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(timestamp),
		AggregateId:  wrapperspb.String(userId.String()),
	})
	if err != nil {
		f.log.Error(err, "failed to create user snapshot", "userId", userId, "timeStamp", timestamp)
		return userOverview
	}
	user, ok := userSnapshot.(*projections.User)
	if !ok {
		f.log.Error(err, "failed to cast user snapshot to user projection", "userId", userId, "timeStamp", timestamp)
		return userOverview
	}
	userOverview.Name = user.Name
	userOverview.Email = user.Email

	overviewFormatter := overviews.NewUserOverviewFormatter(f.EsClient)
	userOverview.Roles, userOverview.Tenants, userOverview.Clusters, err = overviewFormatter.GetRolesDetails(ctx, user, timestamp)
	if err != nil {
		f.log.Error(err, "failed to format roles details", "userId", user, "timeStamp", timestamp)
	}
	userOverview.Details, err = overviewFormatter.GetFormattedDetails(ctx, user, timestamp)
	if err != nil {
		f.log.Error(err, "failed to format overview details", "userId", user.Id, "timeStamp", timestamp)
	}

	return userOverview
}
