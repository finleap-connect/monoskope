package acceptance

import (
	"github.com/cucumber/godog"
)

func iCreateAClusterWithTheDnsAddress(arg1 string) error {
	return godog.ErrPending
}

func theCommandShouldFail() error {
	return godog.ErrPending
}

func thereIsAnEmptyListOfClusters() error {
	return godog.ErrPending
}

func thereAreClustersWithTheDnsAddressesOf(tbl *Table) error {
	return godog.ErrPending
}

func thereShouldBeAClusterWithTheNameInTheListOfClusters(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I create a cluster with the dns address "([^"]*)"$`, iCreateAClusterWithTheDnsAddress)
	ctx.Step(`^the command should fail\.$`, theCommandShouldFail)
	ctx.Step(`^there is an empty list of clusters$`, thereIsAnEmptyListOfClusters)
	ctx.Step(`^there are clusters with the dns addresses of:$`, thereAreClustersWithTheDnsAddressesOf)
	ctx.Step(`^there should be a cluster with the name "([^"]*)" in the list of clusters\.$`, thereShouldBeAClusterWithTheNameInTheListOfClusters)
}
