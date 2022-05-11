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
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test validation rules for cluster mapping messages", func() {
	var cd *commanddata.CreateTenantClusterBindingCommandData
	var ed *eventdata.TenantClusterBindingCreated
	JustBeforeEach(func() {
		cd = NewValidCreateTenantClusterBindingCommandData()
		ed = NewValidTenantClusterBindingCreated()
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

	It("should check for a valid TenantId", func() {
		cd.TenantId = invalidUUID
		ed.TenantId = invalidUUID
		ValidateErrorExpected()
	})

	It("should check for a valid ClusterId", func() {
		cd.ClusterId = invalidUUID
		ed.ClusterId = invalidUUID
		ValidateErrorExpected()
	})
})
