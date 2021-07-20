package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedCSR                     = []byte("This should be a CSR")
	expectedReferencedAggregateId   = uuid.New()
	expectedReferencedAggregateType = aggregates.Cluster
)

var _ = Describe("Pkg/Domain/Aggregates/Certificate", func() {
	It("should handle RequestCertificateCommand correctly", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		inAggId := uuid.New()
		agg := NewCertificateAggregate(inAggId)

		reply, err := newRequestCertificateCommand(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		// This is a create command and should set a new ID, regardless of what was passed in.
		Expect(reply.Id).ToNot(Equal(inAggId))

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

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New())

		ed := es.ToEventDataFromProto(&eventdata.CertificateRequested{
			ReferencedAggregateId:   expectedReferencedAggregateId.String(),
			ReferencedAggregateType: expectedReferencedAggregateType.String(),
			SigningRequest:          expectedCSR,
		})
		esEvent := es.NewEvent(ctx, events.CertificateRequested, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*CertificateAggregate).signingRequest).To(Equal(expectedCSR))
		Expect(agg.(*CertificateAggregate).relatedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(agg.(*CertificateAggregate).relatedAggregateType).To(Equal(expectedReferencedAggregateType))
	})
})

func newRequestCertificateCommand(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := commands.NewRequestCertificateCommand(uuid.New()).(*commands.RequestCertificateCommand)
	Expect(ok).To(BeTrue())

	esCommand.SigningRequest = expectedCSR
	esCommand.ReferencedAggregateId = expectedReferencedAggregateId.String()
	esCommand.ReferencedAggregateType = expectedReferencedAggregateType.String()

	return agg.HandleCommand(ctx, esCommand)
}
