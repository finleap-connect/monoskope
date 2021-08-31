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

package rabbitmq

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var (
	baseTestEnv *test.TestEnv
	env         *TestEnv
)

func TestRabbitmq(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/rabbitmq-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Rabbitmq Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")

	var err error
	baseTestEnv = test.NewTestEnv("messaging-testenv")
	env, err = NewTestEnvWithParent(baseTestEnv)
	Expect(err).ToNot(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(env.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
