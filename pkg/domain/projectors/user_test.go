package projectors

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var _ = Describe("Testing with Ginkgo", func() {
})

var _ = Describe("domain/user_repo", func() {
	ctx := context.Background()
	userId := uuid.New()
	adminUser := &projections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}

	It("can handle events", func() {
		userProjector := NewUserProjector()
		userProjection := userProjector.NewProjection(uuid.New())
		protoEventData := &eventdata.UserCreated{
			Name:  adminUser.Name,
			Email: adminUser.Email,
		}
		eventData := eventsourcing.ToEventDataFromProto(protoEventData)
		userProjection, err := userProjector.Project(context.Background(), eventsourcing.NewEvent(ctx, events.UserCreated, eventData, time.Now().UTC(), aggregates.User, uuid.MustParse(adminUser.Id), 1), userProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(userProjection.Version()).To(Equal(uint64(1)))
	})
})
