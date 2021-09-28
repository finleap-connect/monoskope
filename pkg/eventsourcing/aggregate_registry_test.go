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

package eventsourcing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("aggregate_registry", func() {
	It("can register and unregistered aggregates", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
	})
	It("can't register the same aggregate twice", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })

		defer func() {
			Expect(recover()).To(HaveOccurred())
		}()
		registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
	})
	It("can't create aggregates which are not registered", func() {
		registry := NewAggregateRegistry()
		aggregate, err := registry.CreateAggregate(testAggregateType)
		Expect(err).To(HaveOccurred())
		Expect(aggregate).To(BeNil())
	})
	It("can create aggregates which are registered", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })

		aggregate, err := registry.CreateAggregate(testAggregateType)
		Expect(err).ToNot(HaveOccurred())
		Expect(aggregate).ToNot(BeNil())
	})
})
