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

package events

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
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

func init() {
	for _, eventType := range events.UserEvents {
		_ = event.DefaultEventFormatterRegistry.RegisterEventFormatter(eventType, NewUserEventFormatter)
	}
}

// userEventFormatter EventFormatter implementation for the user-aggregate
type userEventFormatter struct {
	*event.EventFormatterBase
}

// NewUserEventFormatter creates a new event formatter for the user-aggregate
func NewUserEventFormatter(esClient esApi.EventStoreClient) event.EventFormatter {
	return &userEventFormatter{
		EventFormatterBase: &event.EventFormatterBase{FormatterBase: &formatters.FormatterBase{EsClient: esClient}},
	}
}

// GetFormattedDetails formats the user-aggregate-events in a human-readable format
func (f *userEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.UserDeleted:
		return f.getFormattedDetailsUserDeleted(ctx, event)
	case events.UserRoleBindingDeleted:
		return f.getFormattedDetailsUserRoleBindingDeleted(ctx, event)
	}

	ed, err := es.EventData(event.Data).Unmarshal()
	if err != nil {
		return "", err
	}

	switch ed := ed.(type) {
	case *eventdata.UserCreated:
		return f.getFormattedDetailsUserCreated(event, ed)
	case *eventdata.UserUpdated:
		return f.getFormattedDetailsUserUpdated(ctx, event, ed)
	case *eventdata.UserRoleAdded:
		return f.getFormattedDetailsUserRoleAdded(ctx, event, ed)
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *userEventFormatter) getFormattedDetailsUserCreated(event *esApi.Event, eventData *eventdata.UserCreated) (string, error) {
	return fmt.Sprintf("“%s“ created user “%s“", event.Metadata[auth.HeaderAuthEmail], eventData.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserUpdated(ctx context.Context, event *esApi.Event, eventData *eventdata.UserUpdated) (string, error) {
	userSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: timestamppb.New(event.GetTimestamp().AsTime().Add(time.Duration(-1) * time.Microsecond)), // exclude the update event
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}
	oldUser, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("“%s“ updated the User", event.Metadata[auth.HeaderAuthEmail]))
	f.AppendUpdate("Name", eventData.Name, oldUser.Name, &details)
	return details.String(), nil
}

func (f *userEventFormatter) getFormattedDetailsUserRoleAdded(ctx context.Context, event *esApi.Event, eventData *eventdata.UserRoleAdded) (string, error) {
	userSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId:  &wrapperspb.StringValue{Value: eventData.UserId}},
	)
	if err != nil {
		return "", err
	}
	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ assigned the role “%s“ for scope “%s“ to user “%s“",
		event.Metadata[auth.HeaderAuthEmail], eventData.Role, eventData.Scope, user.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	userSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId:  &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}

	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ deleted user “%s“", event.Metadata[auth.HeaderAuthEmail], user.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserRoleBindingDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: event.AggregateId}
	urbSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserRoleBindingProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	urb, ok := urbSnapshot.(*projections.UserRoleBinding)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: urb.UserId}
	userSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ removed the role “%s“ for scope “%s“ from user “%s“",
		event.Metadata[auth.HeaderAuthEmail], urb.Role, urb.Scope, user.Email), nil
}
