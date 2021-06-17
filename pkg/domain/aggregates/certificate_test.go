package aggregates

import (
	"context"

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
	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewCertificateAggregate(uuid.New())

		err = newRequestCertificateCommand(ctx, agg)
		Expect(err).NotTo(HaveOccurred())

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.CertificateRequested))

		data := &eventdata.CertificateRequested{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.SigningRequest).To(Equal(expectedCSR))
		Expect(data.ReferencedAggregateId).To(Equal(expectedReferencedAggregateId.String()))
		Expect(data.ReferencedAggregateType).To(Equal(expectedReferencedAggregateType.String()))
	})
})

func newRequestCertificateCommand(ctx context.Context, agg es.Aggregate) error {
	esCommand, ok := commands.NewRequestCertificateCommand(uuid.New()).(*commands.RequestCertificateCommand)
	Expect(ok).To(BeTrue())

	esCommand.SigningRequest = expectedCSR
	esCommand.ReferencedAggregateId = expectedReferencedAggregateId.String()
	esCommand.ReferencedAggregateType = expectedReferencedAggregateType.String()

	return agg.HandleCommand(ctx, esCommand)
}
