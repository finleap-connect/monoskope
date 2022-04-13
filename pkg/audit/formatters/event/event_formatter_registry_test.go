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

package event

import (
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("event_formatter_registry", func() {
	eventFormatter := func(esApi.EventStoreClient) EventFormatter { return struct{ EventFormatter }{} }
	esClient := struct{ esApi.EventStoreClient }{}
	eventType := es.EventType("TestEventType")

	It("can register event formatter for event type", func() {
		registry := NewEventFormatterRegistry()
		err := registry.RegisterEventFormatter(eventType, eventFormatter)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can't replace registered event formatter for event type", func() {
		registry := NewEventFormatterRegistry()
		err := registry.RegisterEventFormatter(eventType, eventFormatter)
		Expect(err).ToNot(HaveOccurred())
		err = registry.RegisterEventFormatter(eventType, eventFormatter)
		Expect(err).To(HaveOccurred())
	})
	It("can't create event formatters which are not registered", func() {
		registry := NewEventFormatterRegistry()
		eventFormatter, err := registry.CreateEventFormatter(esClient, eventType)
		Expect(err).To(HaveOccurred())
		Expect(eventFormatter).To(BeNil())
	})
	It("can create event formatters which are registered", func() {
		registry := NewEventFormatterRegistry()
		err := registry.RegisterEventFormatter(eventType, eventFormatter)
		Expect(err).ToNot(HaveOccurred())

		aggregate, err := registry.CreateEventFormatter(esClient, eventType)
		Expect(err).ToNot(HaveOccurred())
		Expect(aggregate).ToNot(BeNil())
	})
})
