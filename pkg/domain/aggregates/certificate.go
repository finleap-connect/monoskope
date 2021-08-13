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
	relatedAggregateId   uuid.UUID
	relatedAggregateType es.AggregateType
	signingRequest       []byte
}

// CertificateAggregate creates a new CertificateAggregate
func NewCertificateAggregate() es.Aggregate {
	return &CertificateAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Cluster),
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

		a.relatedAggregateId = id
		a.relatedAggregateType = es.AggregateType(data.GetReferencedAggregateType())
		a.signingRequest = data.GetSigningRequest()
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
