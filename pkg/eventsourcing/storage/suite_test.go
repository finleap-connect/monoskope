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

package storage

import (
	"testing"

	"github.com/finleap-connect/monoskope/internal/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	baseTestEnv *test.TestEnv
	env         *TestEnv
)

func TestEventStoreStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "eventsourcing/storage")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		var err error

		By("bootstrapping test env")
		baseTestEnv = test.NewTestEnv("storage-testenv")
		env, err = NewTestEnvWithParent(baseTestEnv)
		Expect(err).ToNot(HaveOccurred())
		close(done)
	}()

	Eventually(done, 180).Should(BeClosed())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(env.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
