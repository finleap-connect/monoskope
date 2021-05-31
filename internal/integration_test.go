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
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	mdManager.SetUserInformation(&metadata.UserInformation{
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

	userServiceClient := func() domainApi.UserClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewUserClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	tenantServiceClient := func() domainApi.TenantClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewTenantClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	clusterServiceClient := func() domainApi.ClusterClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewClusterClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	It("manage a user", func() {
		userId := uuid.New()
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(userId, commandTypes.CreateUser),
			&cmdData.CreateUserCommandData{Name: "Jane Doe", Email: "jane.doe@monoskope.io"},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())
		Expect(user).ToNot(BeNil())
		Expect(user.GetEmail()).To(Equal("jane.doe@monoskope.io"))
		Expect(user.Id).To(Equal(userId.String()))

		userRoleBindingId := uuid.New()
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(userRoleBindingId, commandTypes.CreateUserRoleBinding),
			&cmdData.CreateUserRoleBindingCommandData{Role: roles.Admin.String(), Scope: scopes.System.String(), UserId: userId.String()},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		// Creating the same rolebinding again should fail
		command.Id = uuid.New().String()
		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).To(HaveOccurred())

		user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())
		Expect(user).ToNot(BeNil())
		Expect(user.Roles[0].Role).To(Equal(roles.Admin.String()))
		Expect(user.Roles[0].Scope).To(Equal(scopes.System.String()))

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())
		Expect(user).ToNot(BeNil())
	})
	It("fail to create a user which already exists", func() {
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(uuid.New(), commandTypes.CreateUser),
			&cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).To(HaveOccurred())
		Expect(errors.TranslateFromGrpcError(err)).To(Equal(errors.ErrUserAlreadyExists))
	})
	It("manage a tenant", func() {
		user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("admin@monoskope.io"))
		Expect(err).ToNot(HaveOccurred())

		tenantId := uuid.New()
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(tenantId, commandTypes.CreateTenant),
			&cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err := tenantServiceClient().GetByName(ctx, wrapperspb.String("Tenant X"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.GetName()).To(Equal("Tenant X"))
		Expect(tenant.GetPrefix()).To(Equal("tx"))
		Expect(tenant.Id).To(Equal(tenantId.String()))

		command, err = cmd.AddCommandData(
			cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
			&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: "DIIIETER"}},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.Metadata.GetLastModifiedBy()).ToNot(BeNil())
		Expect(tenant.Metadata.GetLastModifiedBy().Id).To(Equal(user.Id))

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
		Expect(err).ToNot(HaveOccurred())
		Expect(tenant).ToNot(BeNil())
		Expect(tenant.Metadata.GetDeletedBy()).ToNot(BeNil())
		Expect(tenant.Metadata.GetDeletedBy().GetId()).To(Equal(user.GetId()))
	})
	It("manage a cluster", func() {
		clusterId := uuid.New()
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(clusterId, commandTypes.CreateCluster),
			&cmdData.CreateCluster{Name: "my awesome cluster", Label: "mac", ApiServerAddress: "my.awesome.cluster", ClusterCACertBundle: []byte("This should be a certificate")},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		Expect(err).ToNot(HaveOccurred())

		// Wait to propagate
		time.Sleep(1000 * time.Millisecond)

		cluster, err := clusterServiceClient().GetByName(ctx, wrapperspb.String("my awesome cluster"))
		Expect(err).ToNot(HaveOccurred())
		Expect(cluster).ToNot(BeNil())
		Expect(cluster.GetName()).To(Equal("my awesome cluster"))
		Expect(cluster.GetLabel()).To(Equal("mac"))
		Expect(cluster.GetApiServerAddress()).To(Equal("my.awesome.cluster"))
		Expect(cluster.GetClusterCACert()).To(Equal([]byte("This should be a certificate")))

		// TODO add reactor to test env for bootstrap token creation
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
