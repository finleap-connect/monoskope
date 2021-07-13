package acceptance

import (
	"context"
	"log"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages-go/v10"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal"
	ch "gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
)

type M8TestEnv struct {
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

func ClusterInitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func ClusterInitializeScenario(ctx *godog.ScenarioContext) {
	RegisterFailHandler(Fail)

	ctx.BeforeScenario(func(*godog.Scenario) {
	})

	client, err := NewM8TestEnv()
	if err != nil {
		fail(err)
	}

	ctx.Step(`^I create a cluster with the dns address "([^"]*)"$`, client.iCreateAClusterWithTheDnsAddress)
	ctx.Step(`^I create a cluster with the dns address "([^"]*)" and the name "([^"]*)"$`, client.iCreateAClusterWithTheDnsAddressAndTheName)

	ctx.Step(`^there is an empty list of clusters$`, client.MustEmptyListOfClusters)
	ctx.Step(`^there should be a cluster with the name "([^"]*)" in the list of clusters$`, client.thereShouldBeAClusterWithTheNameInTheListOfClusters)
	ctx.Step(`^there should be a cluster with the dns address "([^"]*)" in the list of clusters$`, client.thereShouldBeAClusterWithTheDnsAddressInTheListOfClusters)

	ctx.Step(`^the command should fail\.$`, client.theCommandShouldFail)
	ctx.Step(`^there are clusters with the dns addresses of:$`, client.mustAllClustersFromList)
	ctx.Step(`^there should be a role binding for the user "([^"]*)" and the role "([^"]*)" for the scope "([^"]*)" and resource "([^"]*)"$`, client.thereShouldBeARoleBindingForTheUserAndTheRoleForTheScopeAndResource)
	ctx.Step(`^there should be a user with the name "([^"]*)" in the list of users$`, client.thereShouldBeAUserWithTheNameInTheListOfUsers)
	ctx.Step(`^my name is "([^"]*)", my email is "([^"]*)" and have a token issued by "([^"]*)"$`, client.myNameIsMyEmailIsAndHaveATokenIssuedBy)
	ctx.Step(`^there are clusters with the names of:$`, client.thereAreClustersWithTheNamesOf)
	ctx.Step(`^there should be a JWT token available for the cluster that is valid for the the name "([^"]*)"$`, client.thereShouldBeAJWTTokenAvailableForTheClusterThatIsValidForTheTheName)
}

func beforeSuite() {
	var err error

	baseTestEnv = test.NewTestEnv("acceptance-stenv")
	testEnv, err = internal.NewTestEnv(baseTestEnv)
	if err != nil {
		fail(err)
	}

}

func afterSuite() {
	err := testEnv.Shutdown()
	if err != nil {
		fail(err)
	}
	err = baseTestEnv.Shutdown()
	if err != nil {
		fail(err)
	}

}

func NewM8TestEnv() (*M8TestEnv, error) {
	var err error
	client := &M8TestEnv{}
	ctx := context.Background()

	client.mdManager, err = metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	client.commandHandlerClient = func() es.CommandHandlerClient {
		addr := testEnv.GetComandHandlerEnv().GetApiAddr()
		_, chClient, err := ch.NewServiceClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	client.userServiceClient = func() domainApi.UserClient {
		addr := testEnv.GetComandHandlerEnv().GetApiAddr()
		_, client, err := queryhandler.NewUserClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	client.tenantServiceClient = func() domainApi.TenantClient {
		addr := testEnv.GetComandHandlerEnv().GetApiAddr()
		_, client, err := queryhandler.NewTenantClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	client.clusterServiceClient = func() domainApi.ClusterClient {
		addr := testEnv.GetComandHandlerEnv().GetApiAddr()
		_, client, err := queryhandler.NewClusterClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	return client, nil
}

func (client *M8TestEnv) iCreateAClusterWithTheDnsAddressAndTheName(arg1, arg2 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) myNameIsMyEmailIsAndHaveATokenIssuedBy(arg1, arg2, arg3 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) mustAllClustersFromList(arg1 *messages.PickleStepArgument_PickleTable) error {
	return godog.ErrPending
}

func (client *M8TestEnv) thereShouldBeAClusterWithTheDnsAddressInTheListOfClusters(arg1 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) thereShouldBeAJWTTokenAvailableForTheClusterThatIsValidForTheTheName(arg1 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) iCreateAClusterWithTheDnsAddress(arg1 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) theCommandShouldFail() error {
	return godog.ErrPending
}

func (client *M8TestEnv) MustEmptyListOfClusters() error {

	return godog.ErrPending
}

func (client *M8TestEnv) thereAreClustersWithTheDnsAddressesOf(tbl *Table) error {
	return godog.ErrPending
}

func (client *M8TestEnv) thereShouldBeAClusterWithTheNameInTheListOfClusters(arg1 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) thereShouldBeARoleBindingForTheUserAndTheRoleForTheScopeAndResource(arg1, arg2, arg3, arg4 string) error {
	return godog.ErrPending
}

func (client *M8TestEnv) thereShouldBeAUserWithTheNameInTheListOfUsers(arg1 string) error {
	return godog.ErrPending
}

// client.mdManager.SetUserInformation(&metadata.UserInformation{
// 	Name:   "admin",
// 	Email:  "admin@monoskope.io",
// 	Issuer: "monoskope",
// })

func fail(err error) {
	log.Fatal(err)
}
