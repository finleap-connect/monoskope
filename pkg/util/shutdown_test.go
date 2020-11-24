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
