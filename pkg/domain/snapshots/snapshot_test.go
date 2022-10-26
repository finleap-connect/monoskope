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

package snapshots

import (
	"context"
	"io"
	"time"

	mock_es "github.com/finleap-connect/monoskope/internal/test/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	meta "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("pkg/domain/snapshots", func() {
	expectedUserId := uuid.New()
	expectedUserName := "Jane Doe"
	expectedUserEmail := "jane.doe@monoskope.io"

	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("can create snapshots", func() {
		metaMgr, err := meta.NewDomainMetadataManager(context.Background())
		Expect(err).NotTo(HaveOccurred())
		metaMgr.SetUserInformation(&meta.UserInformation{
			Id: mock.TestAdminUser.ID(),
		})
		ctx := metaMgr.GetContext()

		timestamp := time.Now()
		maxTimestamp := timestamp.Add(1 * time.Second)
		eventFilter := &esApi.EventFilter{
			AggregateId:  wrapperspb.String(expectedUserId.String()),
			MaxTimestamp: timestamppb.New(maxTimestamp.Add(500 * time.Millisecond)),
		}

		esClient := mock_es.NewMockEventStoreClient(mockCtrl)
		esRetrieveClient := mock_es.NewMockEventStore_RetrieveClient(mockCtrl)
		esClient.EXPECT().Retrieve(ctx, eventFilter).Return(esRetrieveClient, nil)
		esRetrieveClient.EXPECT().Recv().Return(es.NewProtoFromEvent(es.NewEvent(ctx, events.UserCreated, es.ToEventDataFromProto(&eventdata.UserCreated{
			Email:  expectedUserEmail,
			Name:   "not expected user name",
			Source: common.UserSource_INTERNAL,
		}), timestamp, aggregates.User, expectedUserId, 1)), nil)
		esRetrieveClient.EXPECT().Recv().Return(es.NewProtoFromEvent(es.NewEvent(ctx, events.UserUpdated, es.ToEventDataFromProto(&eventdata.UserUpdated{
			Name: expectedUserName,
		}), maxTimestamp, aggregates.User, expectedUserId, 2)), nil)
		esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF)

		userSnapshot := NewSnapshot(esClient, projectors.NewUserProjector())
		user, err := userSnapshot.Create(ctx, eventFilter)
		Expect(err).ToNot(HaveOccurred())
		Expect(user.Name).To(Equal(expectedUserName))
	})
})
