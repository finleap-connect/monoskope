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

type testObserver struct {
	stack []*testProjection
}

func (t *testObserver) expect(p *testProjection) {
	t.stack = append(t.stack, p)
}

func (t *testObserver) Notify(p *testProjection) {
	if len(t.stack) < 1 {
		panic("stack empty")
	}

	n := len(t.stack) - 1
	expected := t.stack[n]
	t.stack = t.stack[:n] // Pop

	if expected != p {
		panic("unexpected")
	}
}

var _ = Describe("repositories/in_memory", func() {
	tp := newTestProjection(uuid.New())
	testReadWrite := func(repo es.Repository[*testProjection]) {
		to := new(testObserver)
		repo.RegisterObserver(to)

		to.expect(tp)
		err := repo.Upsert(context.Background(), tp)
		Expect(err).NotTo(HaveOccurred())

		repo.DeregisterObserver(to)
		err = repo.Upsert(context.Background(), tp)
		Expect(err).NotTo(HaveOccurred())

		projection, err := repo.ById(context.Background(), tp.ID())
		Expect(err).NotTo(HaveOccurred())
		Expect(projection).To(Equal(tp))

		projections, err := repo.All(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(projections).ToNot(BeNil())
		Expect(len(projections)).To(BeNumerically("==", 1))
		Expect(projections[0]).To(Equal(tp))
	}

	It("can read/write projections", func() {
		testReadWrite(NewInMemoryRepository[*testProjection]())
	})

})
