package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// CertificateAggregate is an aggregate for certificates.
type CertificateAggregate struct {
	*DomainAggregateBase
	aggregateManager        es.AggregateStore
	referencedAggregateId   uuid.UUID
	referencedAggregateType es.AggregateType
	signingRequest          []byte
	certificate             []byte
	caCertBundle            []byte
}

// CertificateAggregate creates a new CertificateAggregate
func NewCertificateAggregate() es.Aggregate {
	return &CertificateAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Certificate),
		},
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *CertificateAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	if err := a.Authorize(ctx, cmd, uuid.Nil); err != nil {
		return nil, err
	}

	switch cmd := cmd.(type) {
	case *commands.RequestCertificateCommand:
		if a.Exists() {
			return nil, errors.ErrCertificateAlreadyExists
		}
		ed := es.ToEventDataFromProto(&eventdata.CertificateRequested{
			ReferencedAggregateId:   cmd.GetReferencedAggregateId(),
			ReferencedAggregateType: cmd.GetReferencedAggregateType(),
			SigningRequest:          cmd.GetSigningRequest(),
		})

		_ = a.AppendEvent(ctx, events.CertificateRequested, ed)

		reply := &es.CommandReply{
			Id:      a.ID(),
			Version: a.Version(),
		}
		return reply, nil
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *CertificateAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.CertificateRequested:
		data := &eventdata.CertificateRequested{}
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}

		id, err := uuid.Parse(data.GetReferencedAggregateId())
		if err != nil {
			return err
		}

		a.referencedAggregateId = id
		a.referencedAggregateType = es.AggregateType(data.GetReferencedAggregateType())
		a.signingRequest = data.GetSigningRequest()
	case events.CertificateIssued:
		data := &eventdata.CertificateIssued{}
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.certificate = data.Certificate.GetCertificate()
		a.caCertBundle = data.Certificate.GetCa()
	case events.CertificateRequestIssued:
		// ignored as it does not update the aggregate. TODO: the state of the signing should be tracked in the aggregate, and thus in the projection.
	case events.CertificateIssueingFailed:
		// ignored as it does not update the aggregate. TODO: the state of the signing should be tracked in the aggregate, and thus in the projection.
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
