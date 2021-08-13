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
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedUserName = "the one cluster"
	expectedEmail    = "me@example.com"
)

var _ = Describe("Unit Test for User Aggregate", func() {

	It("should set the data from a command to the resultant event", func() {
		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

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

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewUserAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.UserCreated{
			Name:  expectedUserName,
			Email: expectedEmail,
		})
		esEvent := es.NewEvent(ctx, events.UserCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*UserAggregate).Name).To(Equal(expectedUserName))
		Expect(agg.(*UserAggregate).Email).To(Equal(expectedEmail))

	})

})

func createUser(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateUserCommand(uuid.New()).(*cmd.CreateUserCommand)
	Expect(ok).To(BeTrue())

	esCommand.Name = expectedUserName
	esCommand.Email = expectedEmail

	return agg.HandleCommand(ctx, esCommand)
}
