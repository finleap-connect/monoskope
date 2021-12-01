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
	"net/http"

	"github.com/elimity-com/scim"
	"github.com/finleap-connect/monoskope/pkg/api/domain"
	mockdomain "github.com/finleap-connect/monoskope/test/api/domain"
	"github.com/finleap-connect/monoskope/test/api/eventsourcing"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		When("call GetAll() with Count set to zero in params", func() {
			It("returns the total user count", func() {
				commandHandlerClient := eventsourcing.NewMockCommandHandlerClient(mockCtrl)
				userClient := mockdomain.NewMockUserClient(mockCtrl)
				userHandler := NewUserHandler(commandHandlerClient, userClient)

				request, err := http.NewRequestWithContext(ctx, http.MethodPost, "getall", nil)
				Expect(err).ToNot(HaveOccurred())

				userClient.EXPECT().GetCount(ctx, gomock.Any()).Return(&domain.GetCountResult{Count: 1337}, nil)

				page, err := userHandler.GetAll(request, scim.ListRequestParams{Count: 0})
				Expect(err).ToNot(HaveOccurred())

				Expect(page.TotalResults).To(Equal(1337))
			})
		})
	})

})
