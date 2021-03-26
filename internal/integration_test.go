package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	metadataMgr, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	metadataMgr.SetUserInformation(&metadata.UserInformation{
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})

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

	tenantServiceClient := func() domainApi.TenantServiceClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewTenantServiceClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	It("create a user", func() {
		command, err := cmd.CreateCommand(uuid.New(), commandTypes.CreateUser, &cmdData.CreateUserCommandData{Name: "Jane Doe", Email: "jane.doe@monoskope.io"})
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

		user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())
		Expect(user).ToNot(BeNil())
		Expect(user.GetEmail()).To(Equal("jane.doe@monoskope.io"))
	})
	It("fail to create a user which already exists", func() {
		command, err := cmd.CreateCommand(uuid.New(), commandTypes.CreateUser, &cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).To(HaveOccurred())
		Expect(errors.TranslateFromGrpcError(err)).To(Equal(errors.ErrUserAlreadyExists))
	})
	It("create a tenant", func() {
		user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("admin@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())

		command, err := cmd.CreateCommand(uuid.New(), commandTypes.CreateUserRoleBinding, &cmdData.CreateUserRoleBindingCommandData{
			UserId: user.GetId(),
			Role:   "admin",
			Scope:  "system",
		})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		tenantId := uuid.New()
		command, err = cmd.CreateCommand(tenantId, commandTypes.CreateTenant, &cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err := tenantServiceClient().GetByName(ctx, wrapperspb.String("Tenant X"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.GetName()).To(Equal("Tenant X"))
		Expect(tenant.GetPrefix()).To(Equal("tx"))
		Expect(tenant.Id).To(Equal(tenantId.String()))

		command, err = cmd.CreateCommand(tenantId, commandTypes.UpdateTenant, &cmdData.UpdateTenantCommandData{
			Id:     tenant.GetId(),
			Update: &cmdData.UpdateTenantCommandData_Update{Name: &wrapperspb.StringValue{Value: "DIIIETER"}},
		})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.GetLastModifiedBy()).ToNot(BeNil())
		Expect(tenant.GetLastModifiedBy().Id).To(Equal(user.Id))

		command, err = cmd.CreateCommand(tenantId, commandTypes.DeleteTenant, &cmdData.DeleteTenantCommandData{
			Id: tenant.GetId(),
		})
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(metadataMgr.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.GetDeletedBy()).ToNot(BeNil())
		Expect(tenant.GetDeletedBy().GetId()).To(Equal(user.GetId()))
	})
})

var _ = Describe("PrometheusMetrics", func() {
	It("can scrape event store metrics", func() {
		res, err := http.Get(fmt.Sprintf("http://%s/metrics", testEnv.eventStoreTestEnv.MetricsListener.Addr()))
		Expect(err).ToNot(HaveOccurred())
		defer res.Body.Close()
		Expect(res.StatusCode).To(Equal(200))
	})
})
