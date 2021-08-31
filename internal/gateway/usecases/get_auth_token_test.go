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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	mockRepos "gitlab.figo.systems/platform/monoskope/monoskope/test/domain/repositories"
)

var _ = Describe("GetAuthToken", func() {
	var mockCtrl *gomock.Controller
	ctx := context.Background()
	expectedClusterId := uuid.New()
	expectedClusterName := "testcluster"
	expectedClusterApiServerAddress := "https://somecluster.io"
	expectedUserId := uuid.New()
	expectedUserName := "admin"
	expectedUserEmail := "admin@monoskope.io"
	expectedIssuer := "https://someissuer.io"
	expectedValidity := time.Hour * 1

	jwtTestEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
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

	It("can retrieve openid conf", func() {
		userRepo := mockRepos.NewMockUserRepository(mockCtrl)
		clusterRepo := mockRepos.NewMockClusterRepository(mockCtrl)

		request := &api.ClusterAuthTokenRequest{
			ClusterId: expectedClusterId.String(),
			Role:      string(k8s.DefaultRole),
		}
		result := new(api.ClusterAuthTokenResponse)
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), userRepo, clusterRepo, expectedIssuer, expectedValidity)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:    expectedUserId,
			Name:  expectedUserName,
			Email: expectedUserEmail,
		})

		userProjection := projections.NewUserProjection(expectedUserId).(*projections.User)
		userProjection.Id = expectedUserId.String()
		userProjection.Name = expectedUserName
		userProjection.Email = expectedUserEmail

		clusterProjection := projections.NewClusterProjection(expectedClusterId).(*projections.Cluster)
		clusterProjection.Id = expectedClusterId.String()
		clusterProjection.Name = expectedClusterName
		clusterProjection.ApiServerAddress = expectedClusterApiServerAddress

		ctxWithUser := mdManager.GetContext()
		userRepo.EXPECT().ByUserId(ctxWithUser, expectedUserId).Return(userProjection, nil)
		clusterRepo.EXPECT().ByClusterId(ctxWithUser, expectedClusterId.String()).Return(clusterProjection, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeNil())
		Expect(result.AccessToken).ToNot(BeEmpty())
	})

})
