package usecases

import (
	"context"

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
		uc := NewGetAuthTokenUsecase(request, result, jwtTestEnv.CreateSigner(), userRepo, clusterRepo)

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
		clusterProjection.Label = expectedClusterName
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
