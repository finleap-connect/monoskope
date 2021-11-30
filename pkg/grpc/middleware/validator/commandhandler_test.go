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
	"github.com/finleap-connect/monoskope/pkg/api/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test validation rules for commandhanlder messages", func() {
	Context("Permission Model", func() {
		var pm *domain.PermissionModel
		JustBeforeEach(func() {
			pm = NewValidPermissionModel()
		})

		ValidateErrorExpected := func() {
			err := pm.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := pm.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for valid Roles", func() {
			pm.Roles = []string{invalidRole, invalidRole, invalidRole}
			ValidateErrorExpected()
		})

		It("should check for valid Scopes", func() {
			pm.Roles = []string{invalidScope, invalidScope, invalidScope}
			ValidateErrorExpected()
		})
	})

	Context("Policy", func() {
		var pc *domain.Policy
		JustBeforeEach(func() {
			pc = NewValidPolicy()
		})

		ValidateErrorExpected := func() {
			err := pc.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := pc.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid Role", func() {
			pc.Role = invalidRole
			ValidateErrorExpected()
		})

		It("should check for a valid Scope", func() {
			pc.Scope = invalidScope
			ValidateErrorExpected()
		})
	})
})
