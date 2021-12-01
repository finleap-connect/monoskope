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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Test validation rules for tenant messages", func() {
	Context("Creating Tenant", func() {
		var cd *commanddata.CreateTenantCommandData
		var ed *eventdata.TenantCreated
		JustBeforeEach(func() {
			cd = NewValidCreateTenantCommandData()
			ed = NewValidTenantCreated()
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

		It("should check for a valid Prefix", func() {
			By("not starting with a number", func() {
				cd.Prefix = invalidTenantPrefixStartWithNumber
				ed.Prefix = invalidTenantPrefixStartWithNumber
				ValidateErrorExpected()
			})
			By("not being too longr", func() {
				cd.Prefix = invalidTenantPrefixTooLong
				ed.Prefix = invalidTenantPrefixTooLong
				ValidateErrorExpected()
			})
		})
	})

	Context("Updating Tenant", func() {
		var cd *commanddata.UpdateTenantCommandData
		var ed *eventdata.TenantUpdated
		JustBeforeEach(func() {
			cd = NewValidUpdateTenantCommandData()
			ed = NewValidTenantUpdated()
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
			cd.Name = &wrapperspb.StringValue{Value: invalidDisplayNameTooLong}
			ed.Name = &wrapperspb.StringValue{Value: invalidDisplayNameTooLong}
			ValidateErrorExpected()
		})
	})
})
