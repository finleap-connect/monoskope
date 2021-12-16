package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)


type userEventFormatter struct {
	EventFormatter
	ctx   context.Context
	event *esApi.Event
}

func newUserEventFormatter(eventFormatter EventFormatter, ctx context.Context, event *esApi.Event) *userEventFormatter {
	return &userEventFormatter{EventFormatter: eventFormatter, ctx: ctx, event: event}
}

func (f *userEventFormatter) getFormattedDetails() string {
	switch es.EventType(f.event.Type) {
	case events.UserDeleted: return f.getFormattedDetailsUserDeleted()
	case events.UserRoleBindingDeleted: return f.getFormattedDetailsUserRoleBindingDeleted()
	}
	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}
	switch ed.(type) {
	case *eventdata.UserCreated: return f.getFormattedDetailsUserCreated(ed.(*eventdata.UserCreated))
	case *eventdata.UserRoleAdded: return f.getFormattedDetailsUserRoleAdded(ed.(*eventdata.UserRoleAdded))
	}
	return ""
}

func (f *userEventFormatter) getFormattedDetailsUserCreated(eventData *eventdata.UserCreated) string {
	return fmt.Sprintf("%s created user %s", f.event.Metadata["x-auth-email"], eventData.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserRoleAdded(eventData *eventdata.UserRoleAdded) string {
	user, err := f.getUserById(f.ctx, eventData.UserId)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s assigned the role “%s” for scope “%s” to user %s",
		f.event.Metadata["x-auth-email"], eventData.Role, eventData.Scope, user.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserDeleted() string {
	user, err := f.getUserById(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s deleted %s", f.event.Metadata["x-auth-email"], user.Email)
}

func (f *userEventFormatter) getFormattedDetailsUserRoleBindingDeleted() string {
	urb, err := f.getUserRoleBindingById(f.ctx, f.event.AggregateId)
	if err != nil {
		return ""
	}
	user, err := f.getUserById(f.ctx,urb.UserId)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s removed the role “%s” for scope “%s” from user %s",
		f.event.Metadata["x-auth-email"], urb.Role, urb.Scope, user.Email)
}