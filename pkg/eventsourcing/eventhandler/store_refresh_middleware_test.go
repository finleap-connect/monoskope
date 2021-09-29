// Copyright 2021 Monoskope Authors
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
	mock_eventsourcing "github.com/finleap-connect/monoskope/test/api/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Eventsourcing/Eventhandler/StoreRefreshMiddleware", func() {
	ctx := context.Background()

	Context("The refresh middleware continuesly keeps queryside up-to-date", func() {
		expectedEventType := eventsourcing.EventType("TestEvent")
		expectedAggregateType := eventsourcing.AggregateType("TestAggregate")
		expectedAggregateId := uuid.New()

		When("there are events gone missing", func() {
			It("it queries and applies them", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				var wg sync.WaitGroup
				wg.Add(1)
				defer wg.Wait()

				esClient := mock_eventsourcing.NewMockEventStoreClient(mockCtrl)
				esRetrieveClient := mock_eventsourcing.NewMockEventStore_RetrieveClient(mockCtrl)

				esClient.EXPECT().Retrieve(ctx, gomock.Any()).Return(esRetrieveClient, nil)
				esRetrieveClient.EXPECT().Recv().Return(eventsourcing.NewProtoFromEvent(eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 2)), nil)
				esRetrieveClient.EXPECT().Recv().Return(eventsourcing.NewProtoFromEvent(eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 3)), nil)
				esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).Do(func() {
					esClient.EXPECT().Retrieve(ctx, gomock.Any()).Return(esRetrieveClient, nil).AnyTimes()
					esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF).AnyTimes()
					wg.Done()
				})

				eventHandlerChain := eventsourcing.UseEventHandlerMiddleware(NewLoggingEventHandler(), NewEventStoreRefreshMiddleware(esClient, 300*time.Millisecond))
				err := eventHandlerChain.HandleEvent(ctx, eventsourcing.NewEvent(ctx, expectedEventType, nil, time.Now().UTC(), expectedAggregateType, expectedAggregateId, 1))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
