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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testEnv *TestEnv

func TestCommandHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/k8sauthz")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		defer GinkgoRecover()
		By("bootstrapping test env")
		env, err := NewTestEnv()
		Expect(err).NotTo(HaveOccurred())
		testEnv = env
		close(done)
	}()

	Eventually(done, 60).Should(BeClosed())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(testEnv.Shutdown()).To(Succeed())
})
