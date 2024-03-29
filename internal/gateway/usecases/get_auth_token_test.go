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
	mock_repos "github.com/finleap-connect/monoskope/internal/test/domain/repositories"
	api_projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/util"
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
	expectedIssuer := "https://someissuer.io"
	expectedValidity := map[string]time.Duration{
		"default": time.Hour * 1,
	}

	jwtTestEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
	Expect(err).NotTo(HaveOccurred())
	defer util.PanicOnErrorFunc(jwtTestEnv.Shutdown)

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("can retrieve openid conf", func() {
		clusterAccessRepo := mock_repos.NewMockClusterAccessRepository(mockCtrl)

		request := &api.ClusterAuthTokenRequest{
			ClusterId: expectedClusterId.String(),
			Role:      string(k8s.DefaultRole),
		}
		result := new(api.ClusterAuthTokenResponse)
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), clusterAccessRepo, expectedIssuer, expectedValidity)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:        mock.TestAdminUser.ID(),
			Name:      mock.TestAdminUser.Name,
			Email:     mock.TestAdminUser.Email,
			NotBefore: time.Now().UTC(),
		})

		clusterProjection := projections.NewClusterProjection(expectedClusterId)
		clusterProjection.Id = expectedClusterId.String()
		clusterProjection.Name = expectedClusterName
		clusterProjection.ApiServerAddress = expectedClusterApiServerAddress

		clusterAccessProjection := &api_projections.ClusterAccess{
			Cluster: clusterProjection.Proto(),
			Roles: []string{
				string(k8s.DefaultRole),
			},
		}

		ctxWithUser := mdManager.GetContext()
		clusterAccessRepo.EXPECT().GetClustersAccessibleByUserId(ctxWithUser, mock.TestAdminUser.ID()).Return([]*api_projections.ClusterAccess{clusterAccessProjection}, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeNil())
		Expect(result.AccessToken).ToNot(BeEmpty())
	})
	It("can not retrieve openid conf", func() {
		clusterAccessRepo := mock_repos.NewMockClusterAccessRepository(mockCtrl)

		request := &api.ClusterAuthTokenRequest{
			ClusterId: expectedClusterId.String(),
			Role:      string(k8s.DefaultRole),
		}
		result := new(api.ClusterAuthTokenResponse)
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), clusterAccessRepo, expectedIssuer, expectedValidity)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:        mock.TestAdminUser.ID(),
			Name:      mock.TestAdminUser.Name,
			Email:     mock.TestAdminUser.Email,
			NotBefore: time.Now().UTC(),
		})

		clusterProjection := projections.NewClusterProjection(expectedClusterId)
		clusterProjection.Id = expectedClusterId.String()
		clusterProjection.Name = expectedClusterName
		clusterProjection.ApiServerAddress = expectedClusterApiServerAddress

		ctxWithUser := mdManager.GetContext()
		clusterAccessRepo.EXPECT().GetClustersAccessibleByUserId(ctxWithUser, mock.TestAdminUser.ID()).Return([]*api_projections.ClusterAccess{}, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).To(HaveOccurred())
		Expect(result).ToNot(BeNil())
		Expect(result.AccessToken).To(BeEmpty())
	})
	It("can not retrieve openid conf for admin role", func() {
		clusterAccessRepo := mock_repos.NewMockClusterAccessRepository(mockCtrl)

		request := &api.ClusterAuthTokenRequest{
			ClusterId: expectedClusterId.String(),
			Role:      string(k8s.AdminRole),
		}
		result := new(api.ClusterAuthTokenResponse)
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), clusterAccessRepo, expectedIssuer, expectedValidity)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:        mock.TestAdminUser.ID(),
			Name:      mock.TestAdminUser.Name,
			Email:     mock.TestAdminUser.Email,
			NotBefore: time.Now().UTC(),
		})

		clusterProjection := projections.NewClusterProjection(expectedClusterId)
		clusterProjection.Id = expectedClusterId.String()
		clusterProjection.Name = expectedClusterName
		clusterProjection.ApiServerAddress = expectedClusterApiServerAddress

		clusterAccessProjection := &api_projections.ClusterAccess{
			Cluster: clusterProjection.Proto(),
			Roles: []string{
				string(k8s.DefaultRole),
			},
		}

		ctxWithUser := mdManager.GetContext()
		clusterAccessRepo.EXPECT().GetClustersAccessibleByUserId(ctxWithUser, mock.TestAdminUser.ID()).Return([]*api_projections.ClusterAccess{clusterAccessProjection}, nil)

		err := uc.Run(ctxWithUser)
		Expect(err).To(HaveOccurred())
		Expect(result).ToNot(BeNil())
		Expect(result.AccessToken).To(BeEmpty())
	})
})
