package event_sourcing

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
)

var _ = Describe("handler", func() {
	It("can chain command handler", func() {
		cmd := &testCommand{
			TestCommandData: commands.TestCommandData{
				Test:      "test",
				TestCount: 0,
			},
		}
		handlerChain := ChainCommandHandler(
			&testCommandHandler{val: 1},
			&testCommandHandler{val: 2},
			&testCommandHandler{val: 3},
		)
		err := handlerChain.HandleCommand(context.Background(), cmd)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd.TestCommandData.TestCount).To(BeNumerically("==", 3))
		Expect(cmd.TestCommandData.Test).To(Equal("test123"))
	})
})
