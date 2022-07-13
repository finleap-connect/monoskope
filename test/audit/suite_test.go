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

package audit

import (
	"testing"

	"github.com/finleap-connect/monoskope/internal"
	"github.com/finleap-connect/monoskope/internal/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *internal.TestEnv
)

func TestAuditLogHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "audit log")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		By("bootstrapping test env")
		var err error
		baseTestEnv := test.NewTestEnv("auditlog-testenv")
		testEnv, err = internal.NewTestEnv(baseTestEnv)
		Expect(err).To(Not(HaveOccurred()))
		close(done)
	}()
	Eventually(done, 60).Should(BeClosed())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(testEnv.Shutdown()).To(Not(HaveOccurred()))
})
