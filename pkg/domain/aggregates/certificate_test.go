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
		agg := NewCertificateAggregate()

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
		agg := NewCertificateAggregate()

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
		agg := NewCertificateAggregate()

		esEvent := es.NewEvent(ctx, events.CertificateRequestIssued, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
	})
	It("should handle CertificateIssued event correctly", func() {

		ctx := createSysAdminCtx()
		agg := NewCertificateAggregate()

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
		agg := NewCertificateAggregate()

		esEvent := es.NewEvent(ctx, events.CertificateIssueingFailed, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
	})
})
