package eventhandler

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apies "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	mock_eventsourcing "gitlab.figo.systems/platform/monoskope/monoskope/test/api/eventsourcing"
)

var _ = Describe("package eventhandler", func() {
	var mockCtrl *gomock.Controller
	ctx := context.Background()

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("reactorEventHandler", func() {
		expectedEventType := eventsourcing.EventType("TestEventType")
		expectedAggregateType := eventsourcing.AggregateType("TestAggregateType")
		expectedAggregateId := uuid.New()

		When("Some event occurs", func() {
			It("can handle event without errors", func() {
				testReactor := newTestReactor()

				esClient := mock_eventsourcing.NewMockEventStoreClient(mockCtrl)
				esStoreClient := mock_eventsourcing.NewMockEventStore_StoreClient(mockCtrl)
				esClient.EXPECT().Store(gomock.Any()).Return(esStoreClient, nil)
				esStoreClient.EXPECT().Send(gomock.AssignableToTypeOf(new(apies.Event))).Return(nil)
				esStoreClient.EXPECT().CloseAndRecv().Return(nil, nil)

				event := eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)
				handler := NewReactorEventHandler(esClient, testReactor)

				err := handler.HandleEvent(ctx, event)
				Expect(err).NotTo(HaveOccurred())
				time.Sleep(1000 * time.Millisecond)
			})
			It("does not store events without a valid user ID", func() {
				testReactor := newOtherTestReactor()
				esClient := mock_eventsourcing.NewMockEventStoreClient(mockCtrl)
				esStoreClient := mock_eventsourcing.NewMockEventStore_StoreClient(mockCtrl)

				esClient.EXPECT().Store(ctx).MaxTimes(0)
				esStoreClient.EXPECT().Send(gomock.AssignableToTypeOf(new(apies.Event))).MaxTimes(0)

				event := eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)
				handler := NewReactorEventHandler(esClient, testReactor)

				err := handler.HandleEvent(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

type testReactor struct{}

func newTestReactor() eventsourcing.Reactor {
	return new(testReactor)
}

func (r *testReactor) HandleEvent(ctx context.Context, event eventsourcing.Event, events chan<- eventsourcing.Event) error {
	defer close(events)

	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).NotTo(HaveOccurred())
	userInfo := metadataManager.GetUserInformation()
	userInfo.Id = uuid.NewSHA1(uuid.NameSpaceURL, []byte("clusterbootstrapreactor.monoskope.local"))
	metadataManager.SetUserInformation(userInfo)
	ctx = metadataManager.GetContext()

	events <- eventsourcing.NewEvent(ctx, event.EventType(), nil, time.Now().UTC(), event.AggregateType(), event.AggregateID(), event.AggregateVersion()+1)
	return nil
}

type otherTestReactor struct{}

func newOtherTestReactor() eventsourcing.Reactor {
	return new(otherTestReactor)
}

func (r *otherTestReactor) HandleEvent(ctx context.Context, event eventsourcing.Event, events chan<- eventsourcing.Event) error {
	defer close(events)

	events <- eventsourcing.NewEvent(ctx, event.EventType(), nil, time.Now().UTC(), event.AggregateType(), event.AggregateID(), event.AggregateVersion()+1)
	return nil
}
