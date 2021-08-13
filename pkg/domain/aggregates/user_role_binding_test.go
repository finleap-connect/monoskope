package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedScope    = scopes.Tenant
	expectedRole     = roles.Admin
	expectedResource = uuid.New()
	expectedUserId   = uuid.New()
)

var _ = Describe("Unit Test for UserRoleBinding Aggregate", func() {

	var (
		aggManager = NewTestAggregateManager()
	)

	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

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
		Expect(data.Resource).To(Equal(expectedResource.String()))
		Expect(data.Scope).To(Equal(expectedScope.String()))
		Expect(data.Role).To(Equal(expectedRole.String()))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewUserRoleBindingAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.UserRoleAdded{
			UserId:   expectedUserId.String(),
			Role:     expectedRole.String(),
			Scope:    expectedScope.String(),
			Resource: expectedResource.String(),
		})
		esEvent := es.NewEvent(ctx, events.UserRoleBindingCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*UserRoleBindingAggregate).resource).To(Equal(expectedResource))
		Expect(agg.(*UserRoleBindingAggregate).role).To(Equal(expectedRole))
		Expect(agg.(*UserRoleBindingAggregate).scope).To(Equal(expectedScope))
		Expect(agg.(*UserRoleBindingAggregate).userId).To(Equal(expectedUserId))

	})

})

func createUserRoleBinding(ctx context.Context, agg es.Aggregate, userId uuid.UUID) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateUserRoleBindingCommand(uuid.New()).(*cmd.CreateUserRoleBindingCommand)
	Expect(ok).To(BeTrue())

	esCommand.UserId = userId.String()
	esCommand.Role = expectedRole.String()
	esCommand.Scope = expectedScope.String()
	esCommand.Resource = expectedResource.String()

	return agg.HandleCommand(ctx, esCommand)
}
