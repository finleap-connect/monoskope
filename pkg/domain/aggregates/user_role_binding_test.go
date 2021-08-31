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

package aggregates

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var _ = Describe("Unit Test for UserRoleBinding Aggregate", func() {

	var (
		aggManager = NewTestAggregateManager()
	)

	It("should set the data from a command to the resultant event", func() {

		ctx := createSysAdminCtx()

		// prepare a valid user
		user_agg := NewUserAggregate(aggManager)
		ret, err := createUser(ctx, user_agg)
		Expect(err).NotTo(HaveOccurred())
		user_agg.IncrementVersion()
		aggManager.(*aggregateTestStore).Add(user_agg)
		expectedUserId = ret.Id

		agg := NewUserRoleBindingAggregate(aggManager)

		reply, err := createUserRoleBinding(ctx, agg, ret.Id)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(uuid.Nil))
		Expect(reply.Version).To(Equal(uint64(0)))

		agg.IncrementVersion() // otherwise it will not be validated.
		aggManager.(*aggregateTestStore).Add(agg)

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.UserRoleBindingCreated))

		data := &eventdata.UserRoleAdded{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.UserId).To(Equal(expectedUserId.String()))
		Expect(data.Resource).To(Equal(expectedResourceId.String()))
		Expect(data.Scope).To(Equal(expectedTenantScope.String()))
		Expect(data.Role).To(Equal(expectedAdminRole.String()))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx := createSysAdminCtx()
		agg := NewUserRoleBindingAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.UserRoleAdded{
			UserId:   expectedUserId.String(),
			Role:     expectedAdminRole.String(),
			Scope:    expectedTenantScope.String(),
			Resource: expectedResourceId.String(),
		})
		esEvent := es.NewEvent(ctx, events.UserRoleBindingCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*UserRoleBindingAggregate).resource).To(Equal(expectedResourceId))
		Expect(agg.(*UserRoleBindingAggregate).role).To(Equal(expectedAdminRole))
		Expect(agg.(*UserRoleBindingAggregate).scope).To(Equal(expectedTenantScope))
		Expect(agg.(*UserRoleBindingAggregate).userId).To(Equal(expectedUserId))

	})

})
