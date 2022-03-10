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

package formatters

import (
	"context"
	"fmt"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/errors"
	"github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

type certificateEventFormatter struct {
	*eventformatter.BaseEventFormatter
}

func NewCertificateEventFormatter(esClient esApi.EventStoreClient) *certificateEventFormatter {
	return &certificateEventFormatter{
		BaseEventFormatter: &eventformatter.BaseEventFormatter{EsClient: esClient},
	}
}

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
	return fmt.Sprintf("“%s“ issued a certificate request", event.Metadata["x-auth-email"]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateRequested(event *esApi.Event) (string, error) {
	return fmt.Sprintf("“%s“ requested a certificate", event.Metadata["x-auth-email"]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssued(event *esApi.Event) (string, error) {
	return fmt.Sprintf("“%s“ issued a certificate", event.Metadata["x-auth-email"]), nil
}

func (f *certificateEventFormatter) getFormattedDetailsCertificateIssuingFailed(event *esApi.Event) (string, error) {
	return fmt.Sprintf("certificate request issuing faild for “%s“", event.Metadata["x-auth-email"]), nil
}
