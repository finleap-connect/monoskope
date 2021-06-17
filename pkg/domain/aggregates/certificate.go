package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// CertificateAggregate is an aggregate for certificates.
type CertificateAggregate struct {
	DomainAggregateBase
	relatedAggregateId   uuid.UUID
	relatedAggregateType es.AggregateType
	signingRequest       []byte
	ca                   []byte
	certificate          []byte
	key                  []byte
}

// CertificateAggregate creates a new CertificateAggregate
func NewCertificateAggregate(id uuid.UUID) es.Aggregate {
	return &CertificateAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Cluster, id),
		},
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *CertificateAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case *commands.RequestCertificateCommand:
		ed := es.ToEventDataFromProto(&eventdata.CertificateRequested{
			ReferencedAggregateId:   cmd.GetReferencedAggregateId(),
			ReferencedAggregateType: cmd.GetReferencedAggregateType(),
			SigningRequest:          cmd.GetSigningRequest(),
		})
		_ = a.AppendEvent(ctx, events.CertificateRequested, ed)
		return nil
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
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
