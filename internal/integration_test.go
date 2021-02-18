package internal

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ch "gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	est "gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
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

	It("create a user", func() {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		chConn, chClient, err := ch.NewServiceClient(ctx, chAddr)
		Expect(err).ToNot(HaveOccurred())
		defer chConn.Close()

		esAddr := testEnv.eventStoreTestEnv.GetApiAddr()
		esConn, esClient, err := est.NewEventStoreClient(ctx, esAddr)
		Expect(err).ToNot(HaveOccurred())
		defer esConn.Close()

		metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
		Expect(err).ToNot(HaveOccurred())

		err = metadataMgr.SetUserInformation(&metadata.UserInformation{
			Email:   "admin@monoskope.io",
			Subject: "admin",
			Issuer:  "monoskope",
		})
		Expect(err).ToNot(HaveOccurred())

		command, err := cmd.CreateCommand(commandTypes.CreateUser, &cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"})
		Expect(err).ToNot(HaveOccurred())

		_, err = chClient.Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		stream, err := esClient.Retrieve(ctx, &es.EventFilter{
			AggregateType: wrapperspb.String(aggregateTypes.User.String()),
		})
		Expect(err).ToNot(HaveOccurred())

		event, err := stream.Recv()
		Expect(err).To(BeNil())
		Expect(event).ToNot(BeNil())
	})
})
