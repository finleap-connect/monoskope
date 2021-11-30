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

var _ = Describe("Test validation rules for certificate messages", func() {
	var rc *eventdata.CertificateRequested
	var cr *commanddata.RequestCertificate
	JustBeforeEach(func() {
		rc = NewValidRequestedCertificate()
		cr = NewValidCertificateRequest()
	})

	It("should ensure rules are valid", func() {
		err := rc.Validate()
		Expect(err).NotTo(HaveOccurred())
		err = cr.Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should check for a valid ReferencedAggregateId", func() {
		rc.ReferencedAggregateId = invalidUUID
		cr.ReferencedAggregateId = invalidUUID
		err := rc.Validate()
		Expect(err).To(HaveOccurred())
		err = cr.Validate()
		Expect(err).To(HaveOccurred())
	})

	It("should check for a valid ReferencedAggregateType", func() {
		By("not starting with a number", func() {
			rc.ReferencedAggregateType = invalidAggregateTypeStartWithNumber
			cr.ReferencedAggregateType = invalidAggregateTypeStartWithNumber
			err := rc.Validate()
			Expect(err).To(HaveOccurred())
			err = cr.Validate()
			Expect(err).To(HaveOccurred())
		})
		By("not being too long", func() {
			rc.ReferencedAggregateType = invalidAggregateTypeTooLong
			cr.ReferencedAggregateType = invalidAggregateTypeTooLong
			err := rc.Validate()
			Expect(err).To(HaveOccurred())
			err = cr.Validate()
			Expect(err).To(HaveOccurred())
		})
	})

	It("should check for a valid SigningRequest", func() {
		rc.SigningRequest = invalidCSR
		cr.SigningRequest = invalidCSR
		err := rc.Validate()
		Expect(err).To(HaveOccurred())
		err = cr.Validate()
		Expect(err).To(HaveOccurred())
	})
})
