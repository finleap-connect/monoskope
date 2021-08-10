package eventhandler

import (
	"context"
	"io"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	mock_eventsourcing "gitlab.figo.systems/platform/monoskope/monoskope/test/api/eventsourcing"
)

var _ = Describe("Pkg/Eventsourcing/Eventhandler/StoreRefreshMiddleware", func() {
	var mockCtrl *gomock.Controller
	ctx := context.Background()

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("The refresh middleware continuesly keeps queryside up-to-date", func() {
		expectedEventType := eventsourcing.EventType("TestEvent")
		expectedAggregateType := eventsourcing.AggregateType("TestAggregate")
		expectedAggregateId := uuid.New()

		When("there are no gaps", func() {
			It("does nothing", func() {
				esClient := mock_eventsourcing.NewMockEventStoreClient(mockCtrl)
				esRetrieveClient := mock_eventsourcing.NewMockEventStore_RetrieveClient(mockCtrl)

				esClient.EXPECT().Retrieve(gomock.Any(), gomock.Any(), gomock.Any()).Return(esRetrieveClient, nil).AnyTimes()
				esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()

				middleware := NewEventStoreRefreshMiddleware(esClient, time.Millisecond*100)
				evHandler := middleware(NewLoggingEventHandler())
				event := eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)
				err := evHandler.HandleEvent(ctx, event)
				Expect(err).ToNot(HaveOccurred())
				time.Sleep(time.Millisecond * 150)
			})
		})
		When("there are events gone missing", func() {
			It("it queries and applies them", func() {
				esClient := mock_eventsourcing.NewMockEventStoreClient(mockCtrl)
				esRetrieveClient := mock_eventsourcing.NewMockEventStore_RetrieveClient(mockCtrl)

				esClient.EXPECT().Retrieve(gomock.Any(), gomock.Any(), gomock.Any()).Return(esRetrieveClient, nil).AnyTimes()
				esRetrieveClient.EXPECT().Recv().Return(
					eventsourcing.NewProtoFromEvent(
						eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 2),
					),
					nil,
				)
				esRetrieveClient.EXPECT().Recv().Return(
					eventsourcing.NewProtoFromEvent(
						eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 3),
					),
					nil,
				)
				esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()

				middleware := NewEventStoreRefreshMiddleware(esClient, time.Millisecond*100)
				evHandler := middleware(NewLoggingEventHandler())
				event := eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)
				err := evHandler.HandleEvent(ctx, event)
				Expect(err).ToNot(HaveOccurred())

				time.Sleep(time.Millisecond * 150)
			})
		})
	})
})
