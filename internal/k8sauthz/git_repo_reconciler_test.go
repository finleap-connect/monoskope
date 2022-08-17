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
	"log"
	"os"
	"time"

	mock_repositories "github.com/finleap-connect/monoskope/internal/test/domain/repositories"
	api_projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/go-git/go-git/v5"
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

		It("Reconcile() reconciles the git repo", func() {
			userRepo := mock_repositories.NewMockUserRepository(mockCtrl)
			clusterAccessRepo := mock_repositories.NewMockClusterAccessRepository(mockCtrl)

			// Temp dir to clone the repository
			dir, err := os.MkdirTemp("", "m8-git-repo-reconciler")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir) // clean up

			reconcilerConfig := NewReconcilerConfig(dir, "m8-", map[string]string{
				"admin": "cluster-admin",
			})

			r, err := git.PlainCloneContext(context.Background(), dir, false, &git.CloneOptions{
				URL:               "https://github.com/git-fixtures/basic.git",
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			Expect(err).ToNot(HaveOccurred())

			clusterAccessProjection := &api_projections.ClusterAccess{
				Cluster: &api_projections.Cluster{
					Id:   uuid.NewString(),
					Name: "cluster-a",
				},
				Roles: []string{string(k8s.AdminRole)},
			}

			// expected calls to mocks
			userRepo.EXPECT().AllWith(context.Background(), true).Return([]*projections.User{userA, userB}, nil)
			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserId(context.Background(), userA.ID()).Return([]*api_projections.ClusterAccess{clusterAccessProjection}, nil)
			clusterAccessRepo.EXPECT().GetClustersAccessibleByUserId(context.Background(), userB.ID()).Return([]*api_projections.ClusterAccess{clusterAccessProjection}, nil)

			reconciler := NewGitRepoReconciler(reconcilerConfig, userRepo, clusterAccessRepo, r)
			Expect(reconciler.Reconcile(context.Background())).To(Succeed())
		})
	})
})
