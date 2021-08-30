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

package internal

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	baseTestEnv *test.TestEnv
	testEnv     *TestEnv
)

func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../reports/internal-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	baseTestEnv = test.NewTestEnv("integration-testenv")
	testEnv, err = NewTestEnv(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 120)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(testEnv.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
