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

package validator

import (
	"github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Test validation rules for cluster messages", func() {
	Context("Creating cluster", func() {
		var cd *commanddata.CreateCluster
		JustBeforeEach(func() {
			cd = NewValidCreateCluster()
		})

		ValidateErrorExpected := func() {
			err := cd.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cd.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid Name", func() {
			cd.Name = invalidName
			ValidateErrorExpected()
		})

		It("should check for a valid DisplayName", func() {
			By("not being too long", func() {
				cd.Name = invalidDisplayNameTooLong
				ValidateErrorExpected()
			})

			By("not containing white spaces", func() {
				cd.Name = invalidDisplayNameWhiteSpaces
				ValidateErrorExpected()
			})
		})

		It("should check for a valid ApiServerAddress", func() {
			cd.ApiServerAddress = invalidApiServerAddress
			ValidateErrorExpected()
		})
	})

	Context("Updating cluster", func() {
		var cd *commanddata.UpdateCluster
		JustBeforeEach(func() {
			cd = NewValidUpdateCluster()
		})

		ValidateErrorExpected := func() {
			err := cd.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cd.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid Name", func() {
			By("not being too long", func() {
				cd.Name = wrapperspb.String(invalidDisplayNameTooLong)
				ValidateErrorExpected()
			})

			By("not containing white spaces", func() {
				cd.Name = wrapperspb.String(invalidDisplayNameWhiteSpaces)
				ValidateErrorExpected()
			})
		})

		It("should check for a valid ApiServerAddress", func() {
			cd.ApiServerAddress = wrapperspb.String(invalidApiServerAddress)
			ValidateErrorExpected()
		})
	})
})
