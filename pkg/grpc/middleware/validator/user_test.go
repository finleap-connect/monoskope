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
	"github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test validation rules for user messages", func() {
	Context("Creating User", func() {
		var cd *commanddata.CreateUserCommandData
		var ed *eventdata.UserCreated
		JustBeforeEach(func() {
			cd = NewValidCreateUserCommandData()
			ed = NewValidUserCreated()
		})

		ValidateErrorExpected := func() {
			err := cd.Validate()
			Expect(err).To(HaveOccurred())
			err = ed.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cd.Validate()
			Expect(err).NotTo(HaveOccurred())
			err = ed.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid Name", func() {
			cd.Name = invalidDisplayNameTooLong
			ed.Name = invalidDisplayNameTooLong
			ValidateErrorExpected()
		})

		It("should check for a valid Email", func() {
			cd.Email = invalidEmail
			ed.Email = invalidEmail
			ValidateErrorExpected()
		})
	})

	Context("Creating User Role Binding", func() {
		var cd *commanddata.CreateUserRoleBindingCommandData
		var ed *eventdata.UserRoleAdded
		JustBeforeEach(func() {
			cd = NewValidCreateUserRoleBindingCommandData()
			ed = NewValidUserRoleAdded()
		})

		ValidateErrorExpected := func() {
			err := cd.Validate()
			Expect(err).To(HaveOccurred())
			err = ed.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cd.Validate()
			Expect(err).NotTo(HaveOccurred())
			err = ed.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid UserId", func() {
			cd.UserId = invalidUUID
			ed.UserId = invalidUUID
			ValidateErrorExpected()
		})

		It("should check for a valid Role", func() {
			cd.Role = invalidRole
			ed.Role = invalidRole
			ValidateErrorExpected()
		})

		It("should check for a valid Scope", func() {
			cd.Scope = invalidScope
			ed.Scope = invalidScope
			ValidateErrorExpected()
		})

		It("should check for a valid Resource", func() {
			cd.Resource = invalidUUID
			ed.Resource = invalidUUID
			ValidateErrorExpected()
		})
	})
})
