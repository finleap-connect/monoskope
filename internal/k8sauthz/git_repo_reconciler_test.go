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

package k8sauthz

import (
	"context"
	_ "embed"
	"time"

	mock_repositories "github.com/finleap-connect/monoskope/internal/test/domain/repositories"
	api_projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("internal/k8sauthz", func() {
	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("GitRepoReconciler", func() {
		userA := projections.NewUserProjection(uuid.New())
		userA.Name = "test-a"
		userA.Email = "test-a@monoskope.io"
		userB := projections.NewUserProjection(uuid.New())
		userB.Name = "test-b"
		userB.Email = "test-b@monoskope.io"
		userB.Metadata.Deleted = timestamppb.New(time.Now().UTC())
		userC := projections.NewUserProjection(uuid.New())
		userC.Name = "test-c"
		userC.Email = "test-c@monoskope.io"

		It("Reconcile() reconciles the git repo", func() {
			userRepo := mock_repositories.NewMockUserRepository(mockCtrl)
			clusterAccessRepo := mock_repositories.NewMockClusterAccessRepository(mockCtrl)

			reconcilerConfig := NewReconcilerConfig(testEnv.repoDir, "", "m8-", []*ClusterRoleMapping{
				{
					Scope:       api_projections.ClusterRole_CLUSTER.String(),
					Role:        string(k8s.AdminRole),
					ClusterRole: "cluster-admin",
				},
				{
					Scope:       api_projections.ClusterRole_CLUSTER.String(),
					Role:        string(k8s.OnCallRole),
					ClusterRole: "cluster-oncallee",
				},
			})

			clusterAccessProjectionA := &api_projections.ClusterAccessV2{
				Cluster: &api_projections.Cluster{
					Id:   uuid.NewString(),
					Name: "cluster-a",
				},
				ClusterRoles: []*api_projections.ClusterRole{
					{Scope: api_projections.ClusterRole_CLUSTER, Role: "admin"},
				},
			}
			clusterAccessProjectionB := &api_projections.ClusterAccessV2{
				Cluster: &api_projections.Cluster{
					Id:   uuid.NewString(),
					Name: "cluster-a",
				},
				ClusterRoles: []*api_projections.ClusterRole{
					{Scope: api_projections.ClusterRole_CLUSTER, Role: "default"},
				},
			}
			clusterAccessProjectionC := &api_projections.ClusterAccessV2{
				Cluster: &api_projections.Cluster{
					Id:   uuid.NewString(),
					Name: "cluster-a",
				},
				ClusterRoles: []*api_projections.ClusterRole{
					{Scope: api_projections.ClusterRole_TENANT, Role: "oncall"},
				},
			}

			// expected calls to mocks
			userRepo.EXPECT().AllWith(context.Background(), true).Return([]*projections.User{userA, userB, userC}, nil)
			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserIdV2(context.Background(), userA.ID()).Return([]*api_projections.ClusterAccessV2{clusterAccessProjectionA}, nil)
			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserIdV2(context.Background(), userB.ID()).Return([]*api_projections.ClusterAccessV2{clusterAccessProjectionB}, nil)
			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserIdV2(context.Background(), userC.ID()).Return([]*api_projections.ClusterAccessV2{clusterAccessProjectionC}, nil)

			reconciler := NewGitRepoReconciler(reconcilerConfig, userRepo, clusterAccessRepo, testEnv.gitRepo)
			Expect(reconciler.Reconcile(context.Background())).To(Succeed())

			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserIdV2(context.Background(), userA.ID()).Return([]*api_projections.ClusterAccessV2{clusterAccessProjectionA}, nil)
			Expect(reconciler.ReconcileUser(context.Background(), userA)).To(Succeed())
		})
	})
})
