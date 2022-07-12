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

package gateway

import (
	"testing"

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *TestEnv
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gateway/integration")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		defer ginkgo.GinkgoRecover()

		var err error

		By("bootstrapping test env")
		baseTestEnv := test.NewTestEnv("queryhandler-testenv")
		eventStoreTestEnv, err := eventstore.NewTestEnvWithParent(baseTestEnv)
		Expect(err).To(Not(HaveOccurred()))
		testEnv, err = NewTestEnvWithParent(baseTestEnv, eventStoreTestEnv)
		Expect(err).ToNot(HaveOccurred())
		close(done)
	}()

	Eventually(done, 60).Should(BeClosed())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(testEnv.Shutdown()).To(BeNil())
	Expect(testEnv.eventStoreTestEnv.Shutdown()).To(Not(HaveOccurred()))

	testEnv.GrpcServer.Shutdown()
	testEnv.LocalOIDCProviderServer.Shutdown()

	defer testEnv.ApiListenerAPIServer.Close()
	defer testEnv.ApiListenerOIDCProviderServer.Close()
})
