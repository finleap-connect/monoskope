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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Test validation rules for cluster messages", func() {
	Context("Creating cluster", func() {
		var cd *commanddata.CreateCluster
		var ed *eventdata.ClusterCreated
		var edV2 *eventdata.ClusterCreatedV2
		JustBeforeEach(func() {
			cd = NewValidCreateCluster()
			ed = NewValidClusterCreated()
			edV2 = NewValidClusterCreatedV2()
		})

		ValidateErrorExpected := func() {
			err := cd.Validate()
			Expect(err).To(HaveOccurred())
			err = ed.Validate()
			Expect(err).To(HaveOccurred())
			err = edV2.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := cd.Validate()
			Expect(err).NotTo(HaveOccurred())
			err = ed.Validate()
			Expect(err).NotTo(HaveOccurred())
			err = edV2.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for a valid Name", func() {
			cd.Name = invalidName
			ed.Label = invalidName
			edV2.Name = invalidName
			ValidateErrorExpected()
		})

		It("should check for a valid DisplayName", func() {
			cd.DisplayName = invalidDisplayNameTooLong
			ed.Name = invalidDisplayNameTooLong
			edV2.DisplayName = invalidDisplayNameTooLong
			ValidateErrorExpected()
		})

		It("should check for a valid ApiServerAddress", func() {
			cd.ApiServerAddress = invalidApiServerAddress
			ed.ApiServerAddress = invalidApiServerAddress
			edV2.ApiServerAddress = invalidApiServerAddress
			ValidateErrorExpected()
		})
	})

	Context("Updating cluster", func() {
		var cd *commanddata.UpdateCluster
		var ed *eventdata.ClusterUpdated
		JustBeforeEach(func() {
			cd = NewValidUpdateCluster()
			ed = NewValidClusterUpdated()
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

		It("should check for a valid DisplayName", func() {
			cd.DisplayName = &wrapperspb.StringValue{Value: invalidDisplayNameTooLong}
			ed.DisplayName = invalidDisplayNameTooLong
			ValidateErrorExpected()
		})

		It("should check for a valid ApiServerAddress", func() {
			cd.ApiServerAddress = &wrapperspb.StringValue{Value: invalidApiServerAddress}
			ed.ApiServerAddress = invalidApiServerAddress
			ValidateErrorExpected()
		})
	})
})
