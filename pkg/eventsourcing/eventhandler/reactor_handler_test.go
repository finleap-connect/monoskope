package eventhandler

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apies "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
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
				esClient.EXPECT().Store(ctx).Return(esStoreClient, nil)
				esStoreClient.EXPECT().Send(gomock.AssignableToTypeOf(new(apies.Event))).Return(errors.New("test backoff"))
				esClient.EXPECT().Store(ctx).Return(esStoreClient, nil)
				esStoreClient.EXPECT().Send(gomock.AssignableToTypeOf(new(apies.Event))).Return(nil)
				esStoreClient.EXPECT().CloseAndRecv().Return(nil, nil)

				event := eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)
				handler := NewReactorEventHandler(esClient, testReactor)
				err := handler.HandleEvent(ctx, event)
				Expect(err).NotTo(HaveOccurred())
				handler.Stop()
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
	events <- eventsourcing.NewEvent(ctx, event.EventType(), nil, time.Now().UTC(), event.AggregateType(), event.AggregateID(), event.AggregateVersion()+1)
	return nil
}
