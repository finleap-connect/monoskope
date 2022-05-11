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

package util

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util.shutdown", func() {
	It("can expect", func() {
		shutdown := NewShutdownWaitGroup()
		shutdown.Expect()
		Expect(shutdown.IsExpected()).To(BeTrue())

		shutdown = NewShutdownWaitGroup()
		Expect(shutdown.IsExpected()).To(BeFalse())
	})
	It("can wait for waitgroup to finish", func() {
		shutdown := NewShutdownWaitGroup()
		shutdown.Add(1)
		go func() {
			defer GinkgoRecover()
			for !shutdown.IsExpected() {
				time.Sleep(100 * time.Millisecond)
			}
			shutdown.Done() // Notify workgroup
		}()
		go func() {
			defer GinkgoRecover()
			shutdown.Expect()
		}()
		shutdown.Wait()
		Expect(shutdown.IsExpected()).To(BeTrue())
	})
	It("can wait timeout for waitgroup to finish", func() {
		shutdown := NewShutdownWaitGroup()
		shutdown.Add(1)
		success := shutdown.WaitOrTimeout(1 * time.Millisecond)
		Expect(success).To(BeFalse())
	})
})
