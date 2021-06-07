package reactors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var _ = Describe("package reactors", func() {
	Context("ClusterBootstrapReactor", func() {
		ctx := context.Background()
		aggregateType := aggregates.Cluster

		When("ClusterCreated event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(1)
			eventType := events.ClusterCreated
			eventData := &eventdata.ClusterCreated{
				Name:                "TestCluster",
				Label:               "test-cluster",
				ApiServerAddress:    "https://localhost",
				CaCertificateBundle: []byte("somecabundle"),
			}

			It("generates a new cluster bootstrap token", func() {
				reactor := NewClusterBootstrapReactor(JwtTestEnv.CreateSigner())
				evs, err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(evs)).To(BeNumerically("==", 1))

				event := evs[0]
				Expect(event.EventType()).To(Equal(events.ClusterBootstrapTokenCreated))

				eventDataTokenCreated := &eventdata.ClusterBootstrapTokenCreated{}
				err = event.Data().ToProto(eventDataTokenCreated)
				Expect(err).NotTo(HaveOccurred())
				Expect(eventDataTokenCreated.JWT).To(Not(BeEmpty()))
			})
		})
		When("ClusterCertificateRequested event occurs", func() {
			aggregateId := uuid.New()
			aggregateVersion := uint64(2)
			eventType := events.ClusterCertificateRequested
			eventData := &eventdata.ClusterCertificateRequested{
				CertificateSigningRequest: []byte("somesigningrequest"),
			}

			It("generates a new certificate", func() {
				reactor := NewClusterBootstrapReactor(JwtTestEnv.CreateSigner())
				evs, err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(evs)).To(BeNumerically("==", 1))

				event := evs[0]
				Expect(event.EventType()).To(Equal(events.ClusterOperatorCertificateRequestIssued))
			})
		})
	})
})
