// Copyright 2022 Monoskope Authors
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

package scim

import (
	"github.com/elimity-com/scim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg/scim/User", func() {
	expectedUserEmail := "test.user@monoskope.io"
	expectedDisplayName := "test.user"

	When("calling NewUserAttribute()", func() {
		It("converts resource attributes to a user attribute", func() {
			userAttribute, err := NewUserAttribute(scim.ResourceAttributes{
				"userName":    expectedUserEmail,
				"displayName": expectedDisplayName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(userAttribute.UserName).To(Equal(expectedUserEmail))
			Expect(userAttribute.DisplayName).To(Equal(expectedDisplayName))
		})
	})

	When("calling NewUserAttribute()", func() {
		It("it trims leading and trailing whitespaces", func() {
			userAttribute, err := NewUserAttribute(scim.ResourceAttributes{
				"userName":    " test.user@monoskope.io ",
				"displayName": " test.user ",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(userAttribute.UserName).To(Equal(expectedUserEmail))
			Expect(userAttribute.DisplayName).To(Equal(expectedDisplayName))
		})
	})
})
