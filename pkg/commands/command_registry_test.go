package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("commands/command_registry", func() {
	It("can register and unregister commands", func() {
		err := Registry.RegisterCommand(func() Command { return &TestCommand{} })
		Expect(err).ToNot(HaveOccurred())

		err = Registry.UnregisterCommand(TestCommandType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can't unregister commands which are not registered", func() {
		err := Registry.UnregisterCommand(TestCommandType)
		Expect(err).To(HaveOccurred())
	})
	It("can't register the same command twice", func() {
		err := Registry.RegisterCommand(func() Command { return &TestCommand{} })
		Expect(err).ToNot(HaveOccurred())

		err = Registry.RegisterCommand(func() Command { return &TestCommand{} })
		Expect(err).To(HaveOccurred())
	})
	It("can't create commands which are not registered", func() {
		err := Registry.UnregisterCommand(TestCommandType)
		Expect(err).ToNot(HaveOccurred())

		cmd, err := Registry.CreateCommand(TestCommandType)
		Expect(err).To(HaveOccurred())
		Expect(cmd).To(BeNil())
	})
	It("can create commands which are registered", func() {
		err := Registry.RegisterCommand(func() Command { return &TestCommand{} })
		Expect(err).ToNot(HaveOccurred())

		cmd, err := Registry.CreateCommand(TestCommandType)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd).ToNot(BeNil())
	})
})
