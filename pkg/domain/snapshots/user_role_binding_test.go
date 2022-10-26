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
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	meta "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg/domain/user_role_binding", func() {
	expectedUserId := uuid.New()

	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("can create all user role binding snapshots of a user", func() {
		metaMgr, err := meta.NewDomainMetadataManager(context.Background())
		Expect(err).NotTo(HaveOccurred())
		metaMgr.SetUserInformation(&meta.UserInformation{
			Id: mock.TestAdminUser.ID(),
		})
		ctx := metaMgr.GetContext()

		timestamp := time.Now()
		maxTimestamp := timestamp.Add(1 * time.Second)

		userRolesCount := 0
		esClient := mock_es.NewMockEventStoreClient(mockCtrl)
		esRetrieveClient := mock_es.NewMockEventStore_RetrieveClient(mockCtrl)
		esClient.EXPECT().Retrieve(ctx, gomock.Any()).Return(esRetrieveClient, nil)
		esRetrieveClient.EXPECT().Recv().Return(es.NewProtoFromEvent(es.NewEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(&eventdata.UserRoleAdded{
			UserId:   expectedUserId.String(),
			Role:     string(roles.Admin),
			Resource: uuid.New().String(),
		}), timestamp, aggregates.UserRoleBinding, uuid.New(), 1)), nil)
		userRolesCount++
		esRetrieveClient.EXPECT().Recv().Return(es.NewProtoFromEvent(es.NewEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(&eventdata.UserRoleAdded{
			UserId:   uuid.New().String(),
			Role:     string(roles.K8sOperator),
			Resource: uuid.New().String(),
		}), timestamp, aggregates.UserRoleBinding, uuid.New(), 1)), nil)
		esRetrieveClient.EXPECT().Recv().Return(es.NewProtoFromEvent(es.NewEvent(ctx, events.UserRoleBindingCreated, es.ToEventDataFromProto(&eventdata.UserRoleAdded{
			UserId:   expectedUserId.String(),
			Role:     string(roles.OnCall),
			Resource: uuid.New().String(),
		}), maxTimestamp, aggregates.UserRoleBinding, uuid.New(), 1)), nil)
		userRolesCount++
		esRetrieveClient.EXPECT().Recv().Return(nil, io.EOF)

		userRoleBindingSnapshot := NewUserRoleBindingSnapshot(esClient)
		userRoles := userRoleBindingSnapshot.CreateAll(ctx, expectedUserId, maxTimestamp)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(userRoles)).To(Equal(userRolesCount))
		Expect(userRoles[0].UserId).To(Equal(expectedUserId.String()))
		Expect(userRoles[0].Role).To(Equal(string(roles.Admin)))
		Expect(userRoles[1].Role).To(Equal(string(roles.OnCall)))
	})
})
