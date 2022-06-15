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

package projectors

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type certificateProjector struct {
	*domainProjector
}

func NewCertificateProjector() es.Projector[*projections.Certificate] {
	return &certificateProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *certificateProjector) NewProjection(id uuid.UUID) *projections.Certificate {
	return projections.NewCertificateProjection(id)
}

// Project updates the state of the projection according to the given event.
func (c *certificateProjector) Project(ctx context.Context, event es.Event, cert *projections.Certificate) (*projections.Certificate, error) {
	// Apply the changes for the event.
	switch event.EventType() {
	case events.CertificateRequested:
		data := &eventdata.CertificateRequested{}
		if err := event.Data().ToProto(data); err != nil {
			return cert, err
		}
		cert.ReferencedAggregateId = data.GetReferencedAggregateId()
		cert.AggregateType = data.GetReferencedAggregateType()
		cert.SigningRequest = data.GetSigningRequest()

		if err := c.projectCreated(event, cert.DomainProjection); err != nil {
			return nil, err
		}
	case events.CertificateIssued:
		data := &eventdata.CertificateIssued{}
		if err := event.Data().ToProto(data); err != nil {
			return cert, err
		}
		cert.CaCertBundle = data.GetCertificate().GetCa()
		cert.Certificate.Certificate = data.Certificate.GetCertificate()
		cert.CaCertBundle = data.Certificate.GetCa()
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := c.projectModified(event, cert.DomainProjection); err != nil {
		return nil, err
	}
	cert.IncrementVersion()

	return cert, nil
}
