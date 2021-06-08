package reactors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

var _ = Describe("package reactors", func() {
	Context("ClusterBootstrapReactor", func() {
		ctx := context.Background()
		aggregateType := aggregates.Cluster

		testEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
		Expect(err).NotTo(HaveOccurred())
		defer util.PanicOnError(testEnv.Shutdown())

		reactor := NewClusterBootstrapReactor(testEnv.CreateSigner())

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

			It("emits a ClusterBootstrapTokenCreated event", func() {
				eventChannel := make(chan eventsourcing.Event)
				defer close(eventChannel)

				go func() {
					err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
					Expect(err).NotTo(HaveOccurred())
				}()

				event := <-eventChannel
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

			It("emits a ClusterOperatorCertificateRequestIssued event", func() {
				eventChannel := make(chan eventsourcing.Event)
				defer close(eventChannel)

				go func() {
					err := reactor.HandleEvent(ctx, eventsourcing.NewEvent(ctx, eventType, eventsourcing.ToEventDataFromProto(eventData), time.Now().UTC(), aggregateType, aggregateId, aggregateVersion), eventChannel)
					Expect(err).NotTo(HaveOccurred())
				}()

				event := <-eventChannel
				Expect(event.EventType()).To(Equal(events.ClusterOperatorCertificateRequestIssued))
			})
		})
	})
})
