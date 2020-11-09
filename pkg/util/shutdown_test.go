package util

import (
	"os"
	"syscall"
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
	if os.Getenv("GITLAB_CI") == "" { // Test does not work in pipeline because signal is killing hte process
		It("can wait for signal to finish", func() {
			shutdown := NewShutdownWaitGroup()

			shutdown.RegisterSignalHandler(func() {
				shutdown.Expect()
			})

			shutdown.Add(1)
			go func() {
				defer GinkgoRecover()

				err := syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
				Expect(err).ToNot(HaveOccurred())

				shutdown.Done() // Notify workgroup
			}()

			shutdown.Wait()
			Expect(shutdown.IsExpected()).To(BeTrue())
		})
	}
})
