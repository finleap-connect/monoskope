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

var _ = Describe("Unit Test for User Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {
		ctx := createSysAdminCtx()
		agg := NewUserAggregate(NewTestAggregateManager())

		reply, err := createUser(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(uuid.Nil))
		Expect(reply.Version).To(Equal(uint64(0)))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.UserCreated))

		data := &eventdata.UserCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedUserName))
		Expect(data.Email).To(Equal(expectedEmail))
	})
	It("should apply the data from an event to the aggregate", func() {
		ctx := createSysAdminCtx()
		agg := NewUserAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.UserCreated{
			Name:  expectedUserName,
			Email: expectedEmail,
		})
		createEvent := es.NewEvent(ctx, events.UserCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(createEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*UserAggregate).Name).To(Equal(expectedUserName))
		Expect(agg.(*UserAggregate).Email).To(Equal(expectedEmail))

		deleteEvent := es.NewEvent(ctx, events.UserDeleted, nil, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(deleteEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*UserAggregate).Deleted()).To(BeTrue())
	})
})
