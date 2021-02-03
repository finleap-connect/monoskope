package event_sourcing

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("command_registry", func() {
	It("can register and unregister commands", func() {
		registry := NewCommandRegistry()
		err := registry.RegisterCommand(func() Command { return &testCommand{} })
		Expect(err).ToNot(HaveOccurred())

		err = registry.UnregisterCommand(TestCommandType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can't unregister commands which are not registered", func() {
		registry := NewCommandRegistry()
		err := registry.UnregisterCommand(TestCommandType)
		Expect(err).To(HaveOccurred())
	})
	It("can't register the same command twice", func() {
		registry := NewCommandRegistry()
		err := registry.RegisterCommand(func() Command { return &testCommand{} })
		Expect(err).ToNot(HaveOccurred())

		err = registry.RegisterCommand(func() Command { return &testCommand{} })
		Expect(err).To(HaveOccurred())
	})
	It("can't create commands which are not registered", func() {
		registry := NewCommandRegistry()
		cmd, err := registry.CreateCommand(TestCommandType, nil)
		Expect(err).To(HaveOccurred())
		Expect(cmd).To(BeNil())
	})
	It("can create commands which are registered", func() {
		registry := NewCommandRegistry()
		err := registry.RegisterCommand(func() Command { return &testCommand{} })
		Expect(err).ToNot(HaveOccurred())

		proto := &api.TestCommandData{Test: "Hello world!"}
		any := &anypb.Any{}
		err = any.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		cmd, err := registry.CreateCommand(TestCommandType, any)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd).ToNot(BeNil())

		testCmd, ok := cmd.(*testCommand)
		Expect(ok).To(BeTrue())
		Expect(testCmd).ToNot(BeNil())
		Expect(testCmd.Test).To(Equal("Hello world!"))
	})
	It("can register handlers", func() {
		registry := NewCommandRegistry()
		err := registry.SetHandler(newTestCommandHandler(), TestCommandType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can handle commands", func() {
		registry := NewCommandRegistry()

		err := registry.SetHandler(newTestCommandHandler(), TestCommandType)
		Expect(err).ToNot(HaveOccurred())

		err = registry.HandleCommand(context.Background(), &testCommand{
			TestCommandData: api.TestCommandData{Test: "world!"},
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
