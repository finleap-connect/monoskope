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

package k8sauthzreactor

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("internal/k8sauthzreactor", func() {
	var repoURL = "https://monoskope.io/test.git"
	var repoUser = "testuser"
	var repoPassword = "testpassword"
	var clusters = []uuid.UUID{uuid.New(), uuid.New()}

	Context("GitRepository", func() {
		It("NewGitRepository() creates a new instance with defaults", func() {
			repo := NewGitRepository(repoURL, repoUser, repoPassword)
			Expect(repo).NotTo(BeNil())
			Expect(repo.url).To(Equal(repoURL))
			Expect(repo.username).To(Equal(repoUser))
			Expect(repo.password).To(Equal(repoPassword))
			Expect(repo.allClusters).To(BeTrue())
			Expect(len(repo.GetClusters())).To(BeNumerically("==", 0))
		})
		It("SetClusters() sets clusters and adjusts allClusters", func() {
			repo := NewGitRepository(repoURL, repoUser, repoPassword)
			Expect(repo).NotTo(BeNil())

			// set clusters
			repo.SetClusters(clusters)
			Expect(repo.allClusters).To(BeFalse())
			Expect(repo.GetClusters()).To(Equal(clusters))

			// clear clusters
			repo.SetClusters(nil)
			Expect(repo.allClusters).To(BeTrue())
			Expect(len(repo.GetClusters())).To(BeNumerically("==", 0))
		})
	})
})
