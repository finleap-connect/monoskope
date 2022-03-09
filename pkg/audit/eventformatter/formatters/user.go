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
	"google.golang.org/protobuf/types/known/wrapperspb"
)


type userEventFormatter struct {
	*eventformatter.BaseEventFormatter
}

func NewUserEventFormatter(esClient esApi.EventStoreClient) *userEventFormatter {
	return &userEventFormatter{
		BaseEventFormatter: &eventformatter.BaseEventFormatter{EsClient: esClient},
	}
}

func (f *userEventFormatter) GetFormattedDetails(ctx context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.UserDeleted: return f.getFormattedDetailsUserDeleted(ctx, event)
	case events.UserRoleBindingDeleted: return f.getFormattedDetailsUserRoleBindingDeleted(ctx, event)
	}

	ed, err := f.ToPortoFromEventData(event.Data)
	if err != nil {
		return "", err
	}

	switch ed.(type) {
	case *eventdata.UserCreated: return f.getFormattedDetailsUserCreated(event, ed.(*eventdata.UserCreated))
	case *eventdata.UserRoleAdded: return f.getFormattedDetailsUserRoleAdded(ctx, event, ed.(*eventdata.UserRoleAdded))
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *userEventFormatter) getFormattedDetailsUserCreated(event *esApi.Event, eventData *eventdata.UserCreated) (string, error) {
	return fmt.Sprintf("“%s“ created user “%s“", event.Metadata["x-auth-email"], eventData.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserRoleAdded(ctx context.Context, event *esApi.Event, eventData *eventdata.UserRoleAdded) (string, error) {
	userSnapshot, err := f.GetSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: eventData.UserId}},
	)
	if err != nil {
		return "", err
	}

	// TODO: per aggregate snapshot method?
	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ assigned the role “%s“ for scope “%s“ to user “%s“",
		event.Metadata["x-auth-email"], eventData.Role, eventData.Scope, user.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	userSnapshot, err := f.GetSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: event.AggregateId}},
	)
	if err != nil {
		return "", err
	}

	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ deleted user “%s“", event.Metadata["x-auth-email"], user.Email), nil
}

func (f *userEventFormatter) getFormattedDetailsUserRoleBindingDeleted(ctx context.Context, event *esApi.Event) (string, error) {
	eventFilter := &esApi.EventFilter{MaxTimestamp: event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: event.AggregateId}
	urbSnapshot, err := f.GetSnapshot(ctx, projectors.NewUserRoleBindingProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	urb, ok := urbSnapshot.(*projections.UserRoleBinding)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: urb.UserId}
	userSnapshot, err := f.GetSnapshot(ctx, projectors.NewUserProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	user, ok := userSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	return fmt.Sprintf("“%s“ removed the role “%s“ for scope “%s“ from user “%s“",
		event.Metadata["x-auth-email"], urb.Role, urb.Scope, user.Email), nil
}