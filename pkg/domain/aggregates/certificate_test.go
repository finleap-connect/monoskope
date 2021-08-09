package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
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
	expectedCa                      = []byte("this should be the CA for the issued certificate")
	expectedCertificate             = []byte("this should be the certificate")
)

var _ = Describe("Pkg/Domain/Aggregates/Certificate", func() {
	It("should handle RequestCertificateCommand correctly", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		inAggId := uuid.New()
		agg := NewCertificateAggregate(inAggId, NewTestAggregateManager())

		reply, err := newRequestCertificateCommand(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		// This is a create command and should set a new ID, regardless of what was passed in.
		Expect(reply.Id).ToNot(Equal(inAggId))
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

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New(), NewTestAggregateManager())

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
		Expect(agg.(*CertificateAggregate).referencedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(agg.(*CertificateAggregate).referencedAggregateType).To(Equal(expectedReferencedAggregateType))
	})
	It("should handle CertificateRequestIssuer event correctly", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New(), NewTestAggregateManager())

		esEvent := es.NewEvent(ctx, events.CertificateRequestIssued, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
	})
	It("should handle CertificateIssued event correctly", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New(), NewTestAggregateManager())
		cagg := agg.(*CertificateAggregate)
		cagg.signingRequest = expectedCSR
		cagg.referencedAggregateId = expectedReferencedAggregateId
		cagg.referencedAggregateType = expectedReferencedAggregateType

		ed := es.ToEventDataFromProto(&eventdata.CertificateIssued{
			Certificate: &common.CertificateChain{
				Ca:          expectedCa,
				Certificate: expectedCertificate,
			},
		})

		esEvent := es.NewEvent(ctx, events.CertificateIssued, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(cagg.signingRequest).To(Equal(expectedCSR))
		Expect(cagg.referencedAggregateId).To(Equal(expectedReferencedAggregateId))
		Expect(cagg.referencedAggregateType).To(Equal(expectedReferencedAggregateType))
		Expect(cagg.caCertBundle).To(Equal(expectedCa))
		Expect(cagg.certificate).To(Equal(expectedCertificate))

	})
	It("should handle CertificateIssueingFailed event correctly", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New(), NewTestAggregateManager())

		esEvent := es.NewEvent(ctx, events.CertificateIssueingFailed, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())
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
