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

package scimserver

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/elimity-com/scim"
	scim_errors "github.com/elimity-com/scim/errors"
	"github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	domain_errors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	mockdomain "github.com/finleap-connect/monoskope/test/api/domain"
	"github.com/finleap-connect/monoskope/test/api/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("internal/scimserver/UserHandler", func() {
	Context("querying", func() {
		var mockCtrl *gomock.Controller
		ctx := context.Background()

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		When("calling Get()", func() {
			request, err := http.NewRequestWithContext(ctx, http.MethodPost, "get", nil)
			Expect(err).ToNot(HaveOccurred())

			It("returns the user referenced by id", func() {
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)
				expectedUser := &projections.User{
					Id:    uuid.New().String(),
					Name:  "test.user",
					Email: "test.user@monoskope.io",
					Metadata: &projections.LifecycleMetadata{
						Created:      timestamppb.Now(),
						LastModified: timestamppb.Now(),
					},
				}

				userClient.EXPECT().GetById(ctx, gomock.Any()).Return(expectedUser, nil)

				userResource, err := userHandler.Get(request, expectedUser.Id)
				Expect(err).ToNot(HaveOccurred())
				Expect(userResource.ID).To(Equal(expectedUser.Id))
				Expect(userResource.Attributes["userName"]).To(Equal(expectedUser.Email))
			})
			It("returns an error if there is a problem upstream", func() {
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)

				userClient.EXPECT().GetById(ctx, gomock.Any()).Return(nil, domain_errors.ErrUserNotFound)

				_, err = userHandler.Get(request, "somefakeid")
				Expect(err).To(HaveOccurred())

				scimErr, ok := err.(scim_errors.ScimError)
				Expect(ok).To(BeTrue())

				Expect(scimErr.Status).To(Equal(http.StatusNotFound))
			})
		})

		When("calling GetAll()", func() {
			request, err := http.NewRequestWithContext(ctx, http.MethodPost, "getall", nil)
			Expect(err).ToNot(HaveOccurred())

			It("returns the total user count with count set to zero in params", func() {
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)

				userClient.EXPECT().GetCount(ctx, gomock.Any()).Return(&domain.GetCountResult{Count: 1337}, nil)

				page, err := userHandler.GetAll(request, scim.ListRequestParams{Count: 0})
				Expect(err).ToNot(HaveOccurred())

				Expect(page.TotalResults).To(Equal(1337))
			})
			It("returns an error if there is a problem upstream", func() {
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)

				someError := errors.New("some error")
				userClient.EXPECT().GetCount(ctx, gomock.Any()).Return(nil, someError)

				_, err = userHandler.GetAll(request, scim.ListRequestParams{Count: 0})
				Expect(err).To(HaveOccurred())

				scimErr, ok := err.(scim_errors.ScimError)
				Expect(ok).To(BeTrue())

				Expect(scimErr.Detail).To(Equal(someError.Error()))
			})
			It("returns all users", func() {
				expectedUserA := &projections.User{
					Id:    uuid.New().String(),
					Name:  "test.user.a",
					Email: "test.user.a@monoskope.io",
					Metadata: &projections.LifecycleMetadata{
						Created:      timestamppb.Now(),
						LastModified: timestamppb.Now(),
					},
				}
				expectedUserB := &projections.User{
					Id:    uuid.New().String(),
					Name:  "test.user.b",
					Email: "test.user.b@monoskope.io",
					Metadata: &projections.LifecycleMetadata{
						Created:      timestamppb.Now(),
						LastModified: timestamppb.Now(),
					},
				}
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)

				getAllCient := mockdomain.NewMockUser_GetAllClient(mockCtrl)

				userClient.EXPECT().GetCount(ctx, gomock.Any()).Return(&domain.GetCountResult{Count: 2}, nil)
				userClient.EXPECT().GetAll(ctx, gomock.Any()).Return(getAllCient, nil)
				getAllCient.EXPECT().Recv().Return(expectedUserA, nil)
				getAllCient.EXPECT().Recv().Return(expectedUserB, nil)
				getAllCient.EXPECT().Recv().Return(nil, io.EOF)

				page, err := userHandler.GetAll(request, scim.ListRequestParams{Count: 100})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(page.Resources)).To(Equal(2))
				Expect(page.TotalResults).To(Equal(2))
			})
		})
	})

})
