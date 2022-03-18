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

package aggregates

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
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
func NewCertificateAggregate(aggregateManager es.AggregateStore) es.Aggregate {
	return &CertificateAggregate{
		DomainAggregateBase: &DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Certificate),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *CertificateAggregate) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
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
	default:
		return nil, fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
	return a.DefaultReply(), nil
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
