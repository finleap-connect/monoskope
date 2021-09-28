// Copyright 2021 Monoskope Authors
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

package eventsourcing

import (
	"context"

	cmdApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		inId := uuid.New()
		reply, err := registry.HandleCommand(context.Background(), &testCommand{
			aggregateId:     inId,
			TestCommandData: cmdApi.TestCommandData{Test: "world!"},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(inId))
		Expect(reply.Version).To(Equal(uint64(0)))
	})
})
