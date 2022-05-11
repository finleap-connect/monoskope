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

package projectors

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain/projectors/tenant_cluster_binding", func() {
	ctx := context.Background()
	expectedBindingId := uuid.New()
	expectedTenantId := uuid.New()
	expectedClusterId := uuid.New()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())
	mdManager.SetUserInformation(&metadata.UserInformation{
		Id:    uuid.New(),
		Name:  "admin",
		Email: "admin@monoskope.io",
	})
	ctx = mdManager.GetContext()

	projector := NewTenantClusterBindingProjector()
	projection := projector.NewProjection(uuid.New())

	It("can project event TenantClusterBindingCreated", func() {
		protoEventData := &eventdata.TenantClusterBindingCreated{
			TenantId:  expectedTenantId.String(),
			ClusterId: expectedClusterId.String(),
		}
		event := es.NewEvent(ctx, events.TenantClusterBindingCreated, es.ToEventDataFromProto(protoEventData), time.Now().UTC(), aggregates.TenantClusterBinding, expectedBindingId, 1)

		projection, err := projector.Project(context.Background(), event, projection)
		Expect(err).NotTo(HaveOccurred())
		Expect(projection.Version()).To(Equal(uint64(1)))

		binding, ok := projection.(*projections.TenantClusterBinding)
		Expect(ok).To(BeTrue())

		Expect(binding.GetTenantId()).To(Equal(expectedTenantId.String()))
		Expect(binding.GetClusterId()).To(Equal(expectedClusterId.String()))

		dp := binding.DomainProjection
		Expect(dp.Created).ToNot(BeNil())
		Expect(dp.LastModified).ToNot(BeNil())
		Expect(dp.Deleted).To(BeNil())
	})
	It("can project event TenantClusterBindingDeleted", func() {
		event := es.NewEvent(ctx, events.TenantClusterBindingDeleted, nil, time.Now().UTC(), aggregates.TenantClusterBinding, expectedBindingId, 2)

		projection, err := projector.Project(context.Background(), event, projection)
		Expect(err).NotTo(HaveOccurred())
		Expect(projection.Version()).To(Equal(uint64(2)))

		binding, ok := projection.(*projections.TenantClusterBinding)
		Expect(ok).To(BeTrue())

		dp := binding.DomainProjection
		Expect(dp.Created).ToNot(BeNil())
		Expect(dp.LastModified).ToNot(BeNil())
		Expect(dp.Deleted).ToNot(BeNil())
	})
})
