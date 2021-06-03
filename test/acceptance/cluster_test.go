package acceptance

import (
	"context"
	"log"

	"github.com/cucumber/godog"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal"
	ch "gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
)

type M8Client struct {
	mdManager            *metadata.DomainMetadataManager
	commandHandlerClient func() es.CommandHandlerClient
	userServiceClient    func() domainApi.UserClient
	tenantServiceClient  func() domainApi.TenantClient
	clusterServiceClient func() domainApi.ClusterClient
}

var (
	baseTestEnv *test.TestEnv
	testEnv     *internal.TestEnv
)

/* func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../reports/internal-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "integration", []Reporter{junitReporter})
}
*/
/*
 var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	baseTestEnv = test.NewTestEnv("integration-testenv")
	testEnv, err = NewTestEnv(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 120)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(testEnv.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
*/

func NewM8Client() (*M8Client, error) {
	var err error
	client := &M8Client{}
	ctx := context.Background()

	client.mdManager, err = metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	client.commandHandlerClient = func() es.CommandHandlerClient {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		_, chClient, err := ch.NewServiceClient(ctx, chAddr)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	client.userServiceClient = func() domainApi.UserClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewUserClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	client.tenantServiceClient = func() domainApi.TenantClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewTenantClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	client.clusterServiceClient = func() domainApi.ClusterClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewClusterClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	return client, nil
}

func (client *M8Client) iCreateAClusterWithTheDnsAddress(arg1 string) error {
	return godog.ErrPending
}

func (client *M8Client) theCommandShouldFail() error {
	return godog.ErrPending
}

func (client *M8Client) thereIsAnEmptyListOfClusters() error {
	return godog.ErrPending
}

func (client *M8Client) thereAreClustersWithTheDnsAddressesOf(tbl *Table) error {
	return godog.ErrPending
}

func (client *M8Client) thereShouldBeAClusterWithTheNameInTheListOfClusters(arg1 string) error {
	return godog.ErrPending
}

func (client *M8Client) thereShouldBeARoleBindingForTheUserAndTheRoleForTheScopeAndResource(arg1, arg2, arg3, arg4 string) error {
	return godog.ErrPending
}

func (client *M8Client) thereShouldBeAUserWithTheNameInTheListOfUsers(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	client, err := NewM8Client()
	if err != nil {
		log.Fatal(err)
	}

	ctx.Step(`^I create a cluster with the dns address "([^"]*)"$`, client.iCreateAClusterWithTheDnsAddress)
	ctx.Step(`^the command should fail\.$`, client.theCommandShouldFail)
	ctx.Step(`^there is an empty list of clusters$`, client.thereIsAnEmptyListOfClusters)
	ctx.Step(`^there are clusters with the dns addresses of:$`, client.thereAreClustersWithTheDnsAddressesOf)
	ctx.Step(`^there should be a cluster with the name "([^"]*)" in the list of clusters\.$`, client.thereShouldBeAClusterWithTheNameInTheListOfClusters)
	ctx.Step(`^there should be a cluster with the name "([^"]*)" in the list of clusters$`, client.thereShouldBeAClusterWithTheNameInTheListOfClusters)
	ctx.Step(`^there should be a role binding for the user "([^"]*)" and the role "([^"]*)" for the scope "([^"]*)" and resource "([^"]*)"$`, client.thereShouldBeARoleBindingForTheUserAndTheRoleForTheScopeAndResource)
	ctx.Step(`^there should be a user with the name "([^"]*)" in the list of users$`, client.thereShouldBeAUserWithTheNameInTheListOfUsers)
}

// client.mdManager.SetUserInformation(&metadata.UserInformation{
// 	Name:   "admin",
// 	Email:  "admin@monoskope.io",
// 	Issuer: "monoskope",
// })
