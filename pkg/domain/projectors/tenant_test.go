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

var _ = Describe("process tenant", func() {

	var (
		expectedName   = "the one tenant"
		expectedPrefix = "the prefix"
	)

	ctx := context.Background()
	userId := uuid.New()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())
	mdManager.SetUserInformation(&metadata.UserInformation{
		Id:    userId,
		Name:  "admin",
		Email: "admin@monoskope.io",
	})
	ctx = mdManager.GetContext()

	It("can handle TenantCreated events", func() {
		tenantProjector := NewTenantProjector()
		tenantProjection := tenantProjector.NewProjection(uuid.New())
		prototenantCreatedEventData := &eventdata.TenantCreated{
			Name:   expectedName,
			Prefix: expectedPrefix,
		}
		tenantCreatedEventData := es.ToEventDataFromProto(prototenantCreatedEventData)
		tenantProjection, err := tenantProjector.Project(ctx, es.NewEvent(ctx, events.TenantCreated, tenantCreatedEventData, time.Now().UTC(), aggregates.Tenant, uuid.New(), 1), tenantProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(tenantProjection.Version()).To(Equal(uint64(1)))

		tenant, ok := tenantProjection.(*projections.Tenant)
		Expect(ok).To(BeTrue())

		dp := tenant.DomainProjection

		Expect(tenant.GetName()).To(Equal(expectedName))
		Expect(tenant.GetPrefix()).To(Equal(expectedPrefix))

		Expect(dp.Created).ToNot(BeNil())
	})

})
