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

package usecases

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/test"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/util"
	mockRepos "github.com/finleap-connect/monoskope/test/domain/repositories"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	expectedValidity := map[string]time.Duration{
		"default": time.Hour * 1,
	}

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
		clusterRepo := mockRepos.NewMockClusterRepository(mockCtrl)

		request := &api.ClusterAuthTokenRequest{
			ClusterId: expectedClusterId.String(),
			Role:      string(k8s.DefaultRole),
		}
		result := new(api.ClusterAuthTokenResponse)
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), clusterRepo, expectedIssuer, expectedValidity)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:        expectedUserId,
			Name:      expectedUserName,
			Email:     expectedUserEmail,
			NotBefore: time.Now().UTC(),
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
		clusterRepo.EXPECT().ByClusterId(ctxWithUser, expectedClusterId.String()).Return(clusterProjection, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeNil())
		Expect(result.AccessToken).ToNot(BeEmpty())
	})

})
