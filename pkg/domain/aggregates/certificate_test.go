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
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var _ = Describe("Pkg/Domain/Aggregates/Certificate", func() {
	It("should handle RequestCertificateCommand correctly", func() {
		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate(NewTestAggregateManager())

		reply, err := newRequestCertificateCommand(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		// This is a create command and should set a new ID, regardless of what was passed in.
		Expect(reply.Id).ToNot(Equal(uuid.Nil))
		Expect(reply.Id).ToNot(Equal(expectedReferencedAggregateId))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.CertificateRequested))

		data := &eventdata.CertificateRequested{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.SigningRequest).To(Equal(expectedCSR))
		Expect(data.ReferencedAggregateId).To(Equal(expectedReferencedAggregateId.String()))
		Expect(data.ReferencedAggregateType).To(Equal(expectedReferencedAggregateType.String()))
	})

	It("should handle CertificateRequested event correctly", func() {
		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.CertificateRequested{
			ReferencedAggregateId:   expectedReferencedAggregateId.String(),
			ReferencedAggregateType: expectedReferencedAggregateType.String(),
			SigningRequest:          expectedCSR,
		})
		esEvent := es.NewEvent(ctx, events.CertificateRequested, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*CertificateAggregate).signingRequest).To(Equal(expectedCSR))
		Expect(agg.(*CertificateAggregate).referencedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(agg.(*CertificateAggregate).referencedAggregateType).To(Equal(expectedReferencedAggregateType))
	})
	It("should handle CertificateRequestIssuer event correctly", func() {

		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate(NewTestAggregateManager())

		esEvent := es.NewEvent(ctx, events.CertificateRequestIssued, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
	})
	It("should handle CertificateIssued event correctly", func() {

		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate(NewTestAggregateManager())

		cagg := agg.(*CertificateAggregate)
		cagg.signingRequest = expectedCSR
		cagg.referencedAggregateId = expectedReferencedAggregateId
		cagg.referencedAggregateType = expectedReferencedAggregateType

		ed := es.ToEventDataFromProto(&eventdata.CertificateIssued{
			Certificate: &common.CertificateChain{
				Ca:          expectedClusterCACertBundle,
				Certificate: expectedCertificate,
			},
		})

		esEvent := es.NewEvent(ctx, events.CertificateIssued, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(cagg.signingRequest).To(Equal(expectedCSR))
		Expect(cagg.referencedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(cagg.referencedAggregateType).To(Equal(expectedReferencedAggregateType))
		Expect(cagg.caCertBundle).To(Equal(expectedClusterCACertBundle))
		Expect(cagg.certificate).To(Equal(expectedCertificate))

	})
	It("should handle CertificateIssueingFailed event correctly", func() {

		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate(NewTestAggregateManager())

		esEvent := es.NewEvent(ctx, events.CertificateIssueingFailed, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
	})
})
