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

package validator

import (
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test validation rules for command messages", func() {
	Context("Command", func() {
		var cmd *commands.Command
		JustBeforeEach(func() {
			cmd = NewValidCommand()
		})

		ValidateErrorExpected := func() {
			err := cmd.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cmd.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for valid id", func() {
			cmd.Id = invalidUUID
			ValidateErrorExpected()
		})

		It("should check for valid Type", func() {
			cmd.Type = invalidCommandType
			ValidateErrorExpected()
		})
	})
})
