// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventhandler

import (
	"context"
	"errors"
	"time"

	mock_eventsourcing "github.com/finleap-connect/monoskope/internal/test/api/eventsourcing"
	apies "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Eventsourcing/Eventhandler/ReactorEventHandler", func() {
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
				esStoreClient.EXPECT().Send(gomock.AssignableToTypeOf(new(apies.Event))).Return(errors.New("test backoff"))
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

	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).NotTo(HaveOccurred())
	userInfo := metadataManager.GetUserInformation()
	userInfo.Id = uuid.NewSHA1(uuid.NameSpaceURL, []byte("reactor.monoskope.local"))
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
