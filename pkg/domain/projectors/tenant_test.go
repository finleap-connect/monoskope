package projectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
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
		Id:     userId,
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
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
