package projectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var _ = Describe("domain/user_repo", func() {
	userId := uuid.New()
	adminUser := &projections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}

	It("can handle events", func() {
		userProjector := NewUserProjector()
		userProjection := userProjector.NewProjection()
		protoEventData := &ed.UserCreatedEventData{
			Name:  adminUser.Name,
			Email: adminUser.Email,
		}
		eventData, err := event_sourcing.ToEventDataFromProto(protoEventData)
		Expect(err).NotTo(HaveOccurred())
		userProjection, err = userProjector.Project(context.Background(), event_sourcing.NewEvent(events.UserCreated, eventData, time.Now().UTC(), aggregates.User, uuid.MustParse(adminUser.Id), 1), userProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(userProjection.GetAggregateVersion()).To(Equal(uint64(1)))
	})
})
