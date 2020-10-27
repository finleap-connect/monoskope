package gateway

import (
	"fmt"
	"os"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"github.com/ory/dockertest/v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

var (
	pool         *dockertest.Pool
	dexContainer *dockertest.Resource
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Gateway", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	log := logger.WithName("gatewaysetup")
	path, err := os.Getwd()
	Expect(err).ToNot(HaveOccurred())

	By("bootstrapping gateway test env")
	pool, err = dockertest.NewPool("")
	Expect(err).ToNot(HaveOccurred())

	log.Info("spawn dex container")
	options := &dockertest.RunOptions{
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
		Cmd:        []string{"/usr/local/bin/dex", "serve", "/etc/dex/cfg/config.yaml"},
		Mounts:     []string{fmt.Sprintf("%s/config/dex:/etc/dex/cfg", path)},
	}
	dexContainer, err = pool.RunWithOptions(options)
	Expect(err).ToNot(HaveOccurred())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = pool.Purge(dexContainer)
	Expect(err).ToNot(HaveOccurred())
})
