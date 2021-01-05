package messaging

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var (
	env *MessageBusTestEnv
)

func TestMessageBus(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/messaging-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "messaging", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	defer close(done)

	By("bootstrapping test env")
	env = &MessageBusTestEnv{
		TestEnv: test.SetupGeneralTestEnv("TestMessageBus"),
	}

	err = env.CreateDockerPool()
	Expect(err).ToNot(HaveOccurred())

	// Start rabbitmq
	container, err := env.RunWithOptions(&dockertest.RunOptions{
		Name:       "rabbitmq",
		Repository: "gitlab.figo.systems/platform/dependency_proxy/containers/bitnami/rabbitmq",
		Tag:        "3",
	}, func(config *dc.HostConfig) {
		config.RestartPolicy = dc.AlwaysRestart()
		config.LogConfig = dc.LogConfig{
			Type: "journald",
		}
	})
	Expect(err).ToNot(HaveOccurred())

	// create rabbit conn
	env.amqpURL = fmt.Sprintf("amqp://user:bitnami@127.0.0.1:%s", container.GetPort("5672/tcp"))
	env.Log.Info("Waiting for rabbitmq to warm up...")
	time.Sleep(20 * time.Second)
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = env.Shutdown()
	Expect(err).To(BeNil())
})
