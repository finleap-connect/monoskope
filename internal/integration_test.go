package internal

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	domainCommands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	aggregateTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	It("create a user", func() {
		commandHandlerAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		commandHandlerConn, commandHandlerClient, err := commandhandler.NewServiceClient(ctx, commandHandlerAddr)
		Expect(err).ToNot(HaveOccurred())
		defer commandHandlerConn.Close()

		eventStoreAddr := testEnv.eventStoreTestEnv.GetApiAddr()
		eventStoreConn, eventStoreClient, err := eventstore.NewEventStoreClient(ctx, eventStoreAddr)
		Expect(err).ToNot(HaveOccurred())
		defer eventStoreConn.Close()

		request := &commands.CommandRequest{
			Command: &commands.Command{
				Type: commandTypes.CreateUser.String(),
				Data: &anypb.Any{},
			},
		}
		err = request.Command.Data.MarshalFrom(&domainCommands.CreateUserCommand{Name: "admin", Email: "admin@monoskope.io"})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient.Execute(ctx, request)
		Expect(err).ToNot(HaveOccurred())

		stream, err := eventStoreClient.Retrieve(ctx, &eventsourcing.EventFilter{
			AggregateType: wrapperspb.String(aggregateTypes.User.String()),
		})
		Expect(err).ToNot(HaveOccurred())

		event, err := stream.Recv()
		Expect(err).ToNot(HaveOccurred())
		Expect(event).ToNot(BeNil())
	})
})
