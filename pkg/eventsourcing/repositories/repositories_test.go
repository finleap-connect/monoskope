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

package repositories

import (
	"context"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("repositories/in_memory", func() {
	testProjection := newTestProjection(uuid.New())
	testReadWrite := func(repo es.Repository[es.Projection]) {
		err := repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		err = repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		projection, err := repo.ById(context.Background(), testProjection.ID())
		Expect(err).NotTo(HaveOccurred())
		Expect(projection).To(Equal(testProjection))

		projections, err := repo.All(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(projections).ToNot(BeNil())
		Expect(len(projections)).To(BeNumerically("==", 1))
		Expect(projections[0]).To(Equal(testProjection))
	}
	testRemove := func(repo es.Repository[es.Projection]) {
		err := repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		err = repo.Remove(context.Background(), testProjection.ID())
		Expect(err).NotTo(HaveOccurred())

		projections, err := repo.All(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(projections).ToNot(BeNil())
		Expect(len(projections)).To(BeNumerically("==", 0))
	}

	It("can read/write projections", func() {
		testReadWrite(NewInMemoryRepository[es.Projection]())
	})
	It("can remove projections", func() {
		testRemove(NewInMemoryRepository[es.Projection]())
	})
})
