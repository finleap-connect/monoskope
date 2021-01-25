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

		cmd, err := Registry.CreateCommand(TestCommandType, nil)
		Expect(err).To(HaveOccurred())
		Expect(cmd).To(BeNil())
	})
	It("can create commands which are registered", func() {
		err := Registry.RegisterCommand(func() Command { return &TestCommand{} })
		Expect(err).ToNot(HaveOccurred())

		proto := &api.TestCommandData{Test: "Hello world!"}
		any := &anypb.Any{}
		err = any.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		cmd, err := Registry.CreateCommand(TestCommandType, any)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd).ToNot(BeNil())

		testCmd, ok := cmd.(*TestCommand)
		Expect(ok).To(BeTrue())
		Expect(testCmd).ToNot(BeNil())
		Expect(testCmd.Test).To(Equal("Hello world!"))
	})
	It("can register handlers", func() {
		err := Registry.SetHandler(&TestAggregate{}, TestCommandType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can handle commands", func() {
		err := Registry.HandleCommand(context.Background(), &TestCommand{
			TestCommandData: api.TestCommandData{Test: "world!"},
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
