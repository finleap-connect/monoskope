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

package events

import (
	"context"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	fConsts "github.com/finleap-connect/monoskope/pkg/domain/constants/formatters"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

func init() {
	for _, eventType := range events.CertificateEvents {
		_ = event.DefaultEventFormatterRegistry.RegisterEventFormatter(eventType, NewCertificateEventFormatter)
	}
}

// certificateEventFormatter EventFormatter implementation for the certificate-aggregate
type certificateEventFormatter struct {
	*event.EventFormatterBase
	esClient esApi.EventStoreClient
}

// NewCertificateEventFormatter creates a new event formatter for the certificate-aggregate
func NewCertificateEventFormatter(esClient esApi.EventStoreClient) event.EventFormatter {
	return &certificateEventFormatter{
		&event.EventFormatterBase{}, esClient,
	}
}

// GetFormattedDetails formats the certificate-aggregate-events in a human-readable format
func (f *certificateEventFormatter) GetFormattedDetails(_ context.Context, event *esApi.Event) (string, error) {
	switch es.EventType(event.Type) {
	case events.CertificateRequestIssued:
		return f.getFormattedDetailsCertificateRequestIssued(event)
	case events.CertificateRequested:
		return f.getFormattedDetailsCertificateRequested(event)
	case events.CertificateIssued:
		return f.getFormattedDetailsCertificateIssued(event)
	case events.CertificateIssueingFailed:
		return f.getFormattedDetailsCertificateIssuingFailed(event)
	}

	return "", errors.ErrMissingFormatterImplementationForEventType
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateRequestIssued(event *esApi.Event) (string, error) {
	return fConsts.RequestIssuedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateRequested(event *esApi.Event) (string, error) {
	return fConsts.CertificateRequestedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssued(event *esApi.Event) (string, error) {
	return fConsts.CertificateIssuedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssuingFailed(event *esApi.Event) (string, error) {
	return fConsts.CertificateIssuingFailedDetailsFormat.Sprint(event.Metadata[auth.HeaderAuthEmail]), nil
}
