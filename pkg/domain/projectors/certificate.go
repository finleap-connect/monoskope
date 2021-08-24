package projectors

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	apiProjections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type certificateProjector struct {
	*domainProjector
}

func NewCertificateProjector() es.Projector {
	return &certificateProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *certificateProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewCertificateProjection(id)
}

// Project updates the state of the projection according to the given event.
func (c *certificateProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.Certificate)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.CertificateRequested:
		data := &eventdata.CertificateRequested{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.ReferencedAggregateId = data.GetReferencedAggregateId()
		p.AggregateType = data.GetReferencedAggregateType()
		p.SigningRequest = data.GetSigningRequest()

		if err := c.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.CertificateIssued:
		data := &eventdata.CertificateIssued{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.CaCertBundle = data.GetCertificate().GetCa()
		p.Certificate = &apiProjections.Certificate{
			Certificate:  data.Certificate.GetCertificate(),
			CaCertBundle: data.Certificate.GetCa(),
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := c.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
