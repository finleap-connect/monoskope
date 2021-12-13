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
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Test validation rules for eventsourcing messages", func() {
	Context("Event", func() {
		var e *eventsourcing.Event
		JustBeforeEach(func() {
			e = NewValidEvent()
		})

		ValidateErrorExpected := func() {
			err := e.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := e.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for valid Type", func() {
			e.Type = invalidEventType
			ValidateErrorExpected()
		})

		It("should check for valid AggregateId", func() {
			e.AggregateId = invalidUUID
			ValidateErrorExpected()
		})

		It("should check for valid AggregateType", func() {
			e.AggregateType = invalidAggregateTypeStartWithNumber
			ValidateErrorExpected()
		})
	})

	Context("Event Filter", func() {
		var ef *eventsourcing.EventFilter
		JustBeforeEach(func() {
			ef = NewValidEventFilter()
		})

		ValidateErrorExpected := func() {
			err := ef.Validate()
			Expect(err).To(HaveOccurred())
		}

		It("should ensure rules are valid", func() {
			err := ef.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should check for valid AggregateId", func() {
			ef.AggregateId = &wrapperspb.StringValue{Value: invalidUUID}
			ValidateErrorExpected()
		})

		It("should check for valid AggregateType", func() {
			ef.AggregateType = &wrapperspb.StringValue{Value: invalidAggregateTypeStartWithNumber}
			ValidateErrorExpected()
		})
	})
})
