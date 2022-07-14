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
)

var _ = Describe("Test validation rules for certificate messages", func() {
	var cr *commanddata.RequestCertificate
	JustBeforeEach(func() {
		cr = NewValidCertificateRequest()
	})

	ValidateErrorExpected := func() {
		err := cr.Validate()
		Expect(err).To(HaveOccurred())
	}

	It("should ensure rules are valid", func() {
		err := cr.Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should check for a valid ReferencedAggregateId", func() {
		cr.ReferencedAggregateId = invalidUUID
		ValidateErrorExpected()
	})

	It("should check for a valid ReferencedAggregateType", func() {
		By("not starting with a number", func() {
			cr.ReferencedAggregateType = invalidAggregateTypeStartWithNumber
			ValidateErrorExpected()
		})
		By("not being too long", func() {
			cr.ReferencedAggregateType = invalidAggregateTypeTooLong
			ValidateErrorExpected()
		})
	})

	It("should check for a valid SigningRequest", func() {
		cr.SigningRequest = invalidCSR
		ValidateErrorExpected()
	})
})
