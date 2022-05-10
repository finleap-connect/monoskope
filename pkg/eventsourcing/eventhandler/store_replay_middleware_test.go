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
	"io"
	"sync"
	"time"

	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	mock_eventsourcing_api "github.com/finleap-connect/monoskope/test/api/eventsourcing"
	mock_eventsourcing "github.com/finleap-connect/monoskope/test/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Eventsourcing/Eventhandler/EventStoreReplayMiddleware", func() {
	ctx := context.Background()

	Context("The replay middleware fills gaps in the queryside", func() {
		expectedEventType := eventsourcing.EventType("TestEventXYZ")
		expectedAggregateType := eventsourcing.AggregateType("TestAggregateXYZ")
		expectedAggregateId := uuid.New()

		When("there is a gap", func() {
			It("it queries and applies them", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				var wg sync.WaitGroup
				wg.Add(1)
				defer wg.Wait()

				esClient := mock_eventsourcing_api.NewMockEventStoreClient(mockCtrl)
				esRetrieveClient := mock_eventsourcing_api.NewMockEventStore_RetrieveClient(mockCtrl)
				eventHandler := mock_eventsourcing.NewMockEventHandler(mockCtrl)

				eventHandler.EXPECT().HandleEvent(ctx, gomock.Any()).Return(&ProjectionOutdatedError{ProjectionVersion: 0})
				eventHandler.EXPECT().HandleEvent(ctx, gomock.Any()).Return(nil)
				esClient.EXPECT().Retrieve(ctx, gomock.Any()).Return(esRetrieveClient, nil)
				esRetrieveClient.EXPECT().Recv().Return(eventsourcing.NewProtoFromEvent(eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1)), nil)
				esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).Do(func() {
					esClient.EXPECT().Retrieve(ctx, gomock.Any()).Return(esRetrieveClient, nil).AnyTimes()
					esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()
					wg.Done()
				})

				eventHandlerChain := eventsourcing.UseEventHandlerMiddleware(eventHandler, NewEventStoreReplayMiddleware(esClient))
				err := eventHandlerChain.HandleEvent(ctx, eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 2))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
