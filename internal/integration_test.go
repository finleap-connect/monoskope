package internal

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ch "gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	est "gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	aggregateTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	err = metadataMgr.SetUserInformation(&metadata.UserInformation{
		Email:   "admin@monoskope.io",
		Subject: "admin",
		Issuer:  "monoskope",
	})
	Expect(err).ToNot(HaveOccurred())

	commandHandlerClient := func() es.CommandHandlerClient {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		_, chClient, err := ch.NewServiceClient(ctx, chAddr)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	eventStoreClient := func() es.EventStoreClient {
		esAddr := testEnv.eventStoreTestEnv.GetApiAddr()
		_, esClient, err := est.NewEventStoreClient(ctx, esAddr)
		Expect(err).ToNot(HaveOccurred())
		return esClient
	}

	userServiceClient := func() domainApi.UserServiceClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewUserServiceClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	It("create a user", func() {
		command, err := cmd.CreateCommand(commandTypes.CreateUser, &cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		eventStream, err := eventStoreClient().Retrieve(ctx, &es.EventFilter{
			AggregateType: wrapperspb.String(aggregateTypes.User.String()),
		})
		Expect(err).ToNot(HaveOccurred())

		event, err := eventStream.Recv()
		Expect(err).ToNot(HaveOccurred())
		Expect(event).ToNot(BeNil())

		usClient := userServiceClient()
		user, err := usClient.GetByEmail(ctx, wrapperspb.String("admin@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())
		Expect(user).ToNot(BeNil())
		Expect(user.GetEmail()).To(Equal("admin@monoskope.io"))
	})
})
