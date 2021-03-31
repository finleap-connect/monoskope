package eventsourcing

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmdApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("command_registry", func() {
	It("can register and unregister commands", func() {
		registry := NewCommandRegistry()
		registry.RegisterCommand(func(uuid.UUID) Command { return &testCommand{} })
	})
	It("can't register the same command twice", func() {
		registry := NewCommandRegistry()
		registry.RegisterCommand(func(uuid.UUID) Command { return &testCommand{} })

		defer func() {
			Expect(recover()).To(HaveOccurred())
		}()

		registry.RegisterCommand(func(uuid.UUID) Command { return &testCommand{} })
	})
	It("can't create commands which are not registered", func() {
		registry := NewCommandRegistry()
		cmd, err := registry.CreateCommand(uuid.New(), testCommandType, nil)
		Expect(err).To(HaveOccurred())
		Expect(cmd).To(BeNil())
	})
	It("can create commands which are registered", func() {
		registry := NewCommandRegistry()
		registry.RegisterCommand(func(uuid.UUID) Command { return &testCommand{} })

		proto := &cmdApi.TestCommandData{Test: "Hello world!"}
		any := &anypb.Any{}
		err := any.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		cmd, err := registry.CreateCommand(uuid.New(), testCommandType, any)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd).ToNot(BeNil())

		testCmd, ok := cmd.(*testCommand)
		Expect(ok).To(BeTrue())
		Expect(testCmd).ToNot(BeNil())
		Expect(testCmd.Test).To(Equal("Hello world!"))
	})
	It("can register handlers", func() {
		registry := NewCommandRegistry()
		registry.SetHandler(newTestCommandHandler(), testCommandType)
	})
	It("can handle commands", func() {
		registry := NewCommandRegistry()

		registry.SetHandler(newTestCommandHandler(), testCommandType)

		err := registry.HandleCommand(context.Background(), &testCommand{
			TestCommandData: cmdApi.TestCommandData{Test: "world!"},
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
