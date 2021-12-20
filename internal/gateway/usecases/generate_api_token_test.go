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

package usecases

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/test"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/util"
	mockRepos "github.com/finleap-connect/monoskope/test/domain/repositories"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ = Describe("GenerateAPIToken", func() {
	var mockCtrl *gomock.Controller
	ctx := context.Background()
	expectedUserId := uuid.New()
	expectedUserName := "test-user"
	expectedUserEmail := "test-user@monoskope.io"
	expectedIssuer := "https://someissuer.io"
	expectedValidity := 24 * time.Hour

	jwtTestEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestGenerateAPIToken"))
	Expect(err).NotTo(HaveOccurred())
	defer util.PanicOnError(jwtTestEnv.Shutdown())

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("can retrieve an API token", func() {
		userRepo := mockRepos.NewMockUserRepository(mockCtrl)

		request := &api.APITokenRequest{
			AuthorizationScopes: []api.AuthorizationScope{
				api.AuthorizationScope_WRITE_SCIM,
			},
			User:     &api.APITokenRequest_UserId{UserId: expectedUserId.String()},
			Validity: durationpb.New(expectedValidity),
		}
		response := new(api.APITokenResponse)
		uc := NewGenerateAPITokenUsecase(request, response, jwtTestEnv.CreateSigner(), userRepo, expectedIssuer)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:    expectedUserId,
			Name:  expectedUserName,
			Email: expectedUserEmail,
		})

		userProjection := projections.NewUserProjection(expectedUserId).(*projections.User)
		userProjection.Id = expectedUserId.String()
		userProjection.Name = expectedUserName
		userProjection.Email = expectedUserEmail

		ctxWithUser := mdManager.GetContext()
		userRepo.EXPECT().ByUserId(ctxWithUser, expectedUserId).Return(userProjection, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		Expect(response.AccessToken).ToNot(BeEmpty())
	})

})
