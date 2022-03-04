package eventformatter

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)


type certificateEventFormatter struct {
	EventFormatter
	event *esApi.Event
}

func newCertificateEventFormatter(eventFormatter EventFormatter, event *esApi.Event) *certificateEventFormatter {
	return &certificateEventFormatter{EventFormatter: eventFormatter, event: event}
}

func (f *certificateEventFormatter) getFormattedDetails(_ context.Context) string {
	switch es.EventType(f.event.Type) {
	case events.CertificateRequestIssued: return f.getFormattedDetailsCertificateRequestIssued()
	case events.CertificateIssueingFailed: return f.getFormattedDetailsCertificateIssuingFailed()
	}

	ed, ok := toPortoFromEventData(f.event.Data)
	if !ok {
		return ""
	}

	switch ed.(type) {
	case *eventdata.CertificateRequested: return f.getFormattedDetailsCertificateRequested(ed.(*eventdata.CertificateRequested))
	case *eventdata.CertificateIssued: return f.getFormattedDetailsCertificateIssued(ed.(*eventdata.CertificateIssued))
	}

	return ""
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateRequested(_ *eventdata.CertificateRequested) string {
	return fmt.Sprintf("“%s“ requested a certificate", f.event.Metadata["x-auth-email"])
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateRequestIssued() string {
	return fmt.Sprintf("“%s“ issued a certificate request", f.event.Metadata["x-auth-email"])
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssued(_ *eventdata.CertificateIssued) string {
	return fmt.Sprintf("“%s“ issued a certificate", f.event.Metadata["x-auth-email"])
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssuingFailed() string {
	return fmt.Sprintf("certificate request issuing faild for “%s“", f.event.Metadata["x-auth-email"])
}