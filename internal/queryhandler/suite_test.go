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

package queryhandler

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *TestEnv
)

func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/queryhandler-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "queryhandler/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	baseTestEnv := test.NewTestEnv("queryhandler-testenv")
	eventStoreTestEnv, err := eventstore.NewTestEnvWithParent(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))

	testEnv, err = NewTestEnvWithParent(baseTestEnv, eventStoreTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = testEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))

	err = testEnv.eventStoreTestEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))
})
