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
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	_ "github.com/finleap-connect/monoskope/pkg/domain/formatters/events"
	"github.com/finleap-connect/monoskope/pkg/domain/formatters/overviews"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"github.com/finleap-connect/monoskope/pkg/domain/snapshots"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// AuditFormatter is the interface definition for the formatter used by the auditLogServer
type AuditFormatter interface {
	NewHumanReadableEvent(context.Context, *esApi.Event) *audit.HumanReadableEvent
	NewUserOverview(context.Context, uuid.UUID, time.Time) *audit.UserOverview
}

// auditFormatter is the implementation of AuditFormatter used by auditLogServer
type auditFormatter struct {
	log        logger.Logger
	efRegistry event.EventFormatterRegistry
	esClient   esApi.EventStoreClient
}

// NewAuditFormatter creates an auditFormatter
func NewAuditFormatter(esClient esApi.EventStoreClient, efRegistry event.EventFormatterRegistry) *auditFormatter {
	return &auditFormatter{
		logger.WithName("audit-formatter"), efRegistry, esClient,
	}
}

// NewHumanReadableEvent creates a HumanReadableEvent of a given event
func (f *auditFormatter) NewHumanReadableEvent(ctx context.Context, event *esApi.Event) *audit.HumanReadableEvent {
	humanReadableEvent := &audit.HumanReadableEvent{
		Timestamp: event.Timestamp,
		Issuer:    event.Metadata[auth.HeaderAuthEmail],
		IssuerId:  event.Metadata[auth.HeaderAuthId],
		EventType: event.Type,
	}
	eventFormatter, err := f.efRegistry.CreateEventFormatter(f.esClient, es.EventType(event.Type))
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
	overviewFormatter := overviews.NewUserOverviewFormatter(f.esClient)
	userSnapshot := snapshots.NewSnapshot(f.esClient, projectors.NewUserProjector())
	userRoleBindingSnapshot := snapshots.NewUserRoleBindingSnapshot(f.esClient)

	user, err := userSnapshot.Create(ctx, &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(timestamp),
		AggregateId:  wrapperspb.String(userId.String()),
	})
	if err != nil {
		f.log.Error(err, "failed to create user snapshot", "userId", userId, "timeStamp", timestamp)
		return userOverview
	}
	for _, role := range userRoleBindingSnapshot.CreateAll(ctx, userId, timestamp) {
		user.Roles = append(user.Roles, role.Proto())
	}

	userOverview.Name = user.Name
	userOverview.Email = user.Email
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
