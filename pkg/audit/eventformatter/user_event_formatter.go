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
	"google.golang.org/protobuf/types/known/wrapperspb"
)


type userEventFormatter struct {
	EventFormatter
	event *esApi.Event
}

func newUserEventFormatter(eventFormatter EventFormatter, event *esApi.Event) *userEventFormatter {
	return &userEventFormatter{EventFormatter: eventFormatter, event: event}
}

func (f *userEventFormatter) getFormattedDetails(ctx context.Context) string {
	switch es.EventType(f.event.Type) {
	case events.UserDeleted: return f.getFormattedDetailsUserDeleted(ctx)
	case events.UserRoleBindingDeleted: return f.getFormattedDetailsUserRoleBindingDeleted(ctx)
	}

	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}

	switch ed.(type) {
	case *eventdata.UserCreated: return f.getFormattedDetailsUserCreated(ed.(*eventdata.UserCreated))
	case *eventdata.UserRoleAdded: return f.getFormattedDetailsUserRoleAdded(ctx, ed.(*eventdata.UserRoleAdded))
	}

	return ""
}

func (f *userEventFormatter) getFormattedDetailsUserCreated(eventData *eventdata.UserCreated) string {
	return fmt.Sprintf("“%s“ created user “%s“", f.event.Metadata["x-auth-email"], eventData.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserRoleAdded(ctx context.Context, eventData *eventdata.UserRoleAdded) string {
	userSnapshot, err := f.getSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: f.event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: eventData.UserId}},
	)
	user, ok := userSnapshot.(*projections.User)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ assigned the role “%s“ for scope “%s“ to user “%s“",
		f.event.Metadata["x-auth-email"], eventData.Role, eventData.Scope, user.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserDeleted(ctx context.Context) string {
	userSnapshot, err := f.getSnapshot(ctx, projectors.NewUserProjector(), &esApi.EventFilter{
		MaxTimestamp: f.event.GetTimestamp(),
		AggregateId: &wrapperspb.StringValue{Value: f.event.AggregateId}},
	)
	user, ok := userSnapshot.(*projections.User)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ deleted user “%s“", f.event.Metadata["x-auth-email"], user.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserRoleBindingDeleted(ctx context.Context) string {
	eventFilter := &esApi.EventFilter{MaxTimestamp: f.event.GetTimestamp()}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: f.event.AggregateId}
	urbSnapshot, err := f.getSnapshot(ctx, projectors.NewUserRoleBindingProjector(), eventFilter)
	urb, ok := urbSnapshot.(*projections.UserRoleBinding)
	if err != nil || !ok {
		return ""
	}
	eventFilter.AggregateId = &wrapperspb.StringValue{Value: urb.UserId}
	userSnapshot, err := f.getSnapshot(ctx, projectors.NewUserProjector(), eventFilter)
	user, ok := userSnapshot.(*projections.User)
	if err != nil || !ok {
		return ""
	}

	return fmt.Sprintf("“%s“ removed the role “%s“ for scope “%s“ from user “%s“",
		f.event.Metadata["x-auth-email"], urb.Role, urb.Scope, user.Email)
}