package internal

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("integration", func() {
	It("create a user", func() {
		testEnv.commandHandlerTestEnv.GetApiAddr()
	})
})
