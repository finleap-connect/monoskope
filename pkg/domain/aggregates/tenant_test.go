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

package aggregates

import (
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	expectedTenantName = "the one tenant"
	expectedPrefix     = "tenant-one"
)

var _ = Describe("Unit Test for the Tenant Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {
		ctx := createSysAdminCtx()
		inID := uuid.Nil
		agg := NewTenantAggregate(NewTestAggregateManager())

		reply, err := createTenant(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(inID))
		Expect(reply.Version).To(Equal(uint64(0)))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.TenantCreated))
		Expect(event.AggregateID()).ToNot(Equal(inID))

		data := &eventdata.TenantCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedTenantName))

	})

	It("should apply the data from an event to the aggregate", func() {
		ctx := createSysAdminCtx()
		agg := NewTenantAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.TenantCreated{
			Name:   expectedTenantName,
			Prefix: expectedPrefix,
		})
		esEvent := es.NewEvent(ctx, events.TenantCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*TenantAggregate).name).To(Equal(expectedTenantName))
		Expect(agg.(*TenantAggregate).prefix).To(Equal(expectedPrefix))

	})
})
