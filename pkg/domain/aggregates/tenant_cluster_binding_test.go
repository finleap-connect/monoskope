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
	esMock "github.com/finleap-connect/monoskope/internal/test/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnitTest: TenantClusterBindingAggregate", func() {
	var mockCtrl *gomock.Controller
	ctx := createSysAdminCtx()

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should handle the CreateTenantClusterBindingCommand correctly", func() {
		aggManager := esMock.NewMockAggregateStore(mockCtrl)

		// Setup aggregates
		tenant := NewTenantAggregate(aggManager)
		cluster := NewClusterAggregate(aggManager)
		binding := NewTenantClusterBindingAggregate(aggManager)

		// Let them be valid
		tenant.IncrementVersion()
		cluster.IncrementVersion()

		// Build create command for binding
		createCommand := commands.NewCreateTenantClusterBindingCommand(uuid.Nil).(*commands.CreateTenantClusterBindingCommand)
		createCommand.TenantId = tenant.ID().String()
		createCommand.ClusterId = cluster.ID().String()

		// Define expected calls to mock
		aggManager.EXPECT().Get(ctx, aggregates.Tenant, tenant.ID()).Return(tenant, nil)
		aggManager.EXPECT().Get(ctx, aggregates.Cluster, cluster.ID()).Return(cluster, nil)
		aggManager.EXPECT().All(ctx, aggregates.TenantClusterBinding).Return(make([]es.Aggregate, 0), nil)

		// Let aggregate handle the create command
		reply, err := binding.HandleCommand(ctx, createCommand)
		Expect(err).ToNot(HaveOccurred())
		Expect(reply).ToNot(BeNil())

		// Validate emitted event(s)
		event := binding.UncommittedEvents()[0]
		Expect(event.EventType()).To(Equal(events.TenantClusterBindingCreated))

		// Validate event data of emitted event(s)
		data := &eventdata.TenantClusterBindingCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(data.TenantId).To(Equal(tenant.ID().String()))
		Expect(data.ClusterId).To(Equal(cluster.ID().String()))

		err = binding.ApplyEvent(event)
		Expect(err).NotTo(HaveOccurred())
		binding.IncrementVersion()

		deleteCommand := commands.NewDeleteTenantClusterBindingCommand(reply.Id)
		reply, err = binding.HandleCommand(ctx, deleteCommand)
		Expect(err).ToNot(HaveOccurred())
		Expect(reply).ToNot(BeNil())

		// Validate emitted event(s)
		event = binding.UncommittedEvents()[0]
		Expect(event.EventType()).To(Equal(events.TenantClusterBindingDeleted))

		err = binding.ApplyEvent(event)
		Expect(err).NotTo(HaveOccurred())
		binding.IncrementVersion()

		Expect(binding.Deleted()).To(BeTrue())
		Expect(binding.Version()).To(BeNumerically("==", 2))
	})
})
