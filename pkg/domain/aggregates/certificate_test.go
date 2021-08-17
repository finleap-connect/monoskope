package aggregates

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		Expect(agg.(*CertificateAggregate).relatedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(agg.(*CertificateAggregate).relatedAggregateType).To(Equal(expectedReferencedAggregateType))
	})
})
