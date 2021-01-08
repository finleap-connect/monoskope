package messaging

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var (
	env *messageBusTestEnv
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
	env = &messageBusTestEnv{
		TestEnv: test.NewTestEnv("TestMessageBus"),
	}

	err = env.CreateDockerPool()
	Expect(err).ToNot(HaveOccurred())

	// Start rabbitmq
	container, err := env.Run(&dockertest.RunOptions{
		Name:       "rabbitmq",
		Repository: "gitlab.figo.systems/platform/dependency_proxy/containers/bitnami/rabbitmq",
		Tag:        "3",
	})
	Expect(err).ToNot(HaveOccurred())

	// Build connection string
	env.amqpURL = fmt.Sprintf("amqp://user:bitnami@127.0.0.1:%s", container.GetPort("5672/tcp"))

	// Wait for rabbitmq to start
	for i := 40; i > 0; i-- {
		env.Log.Info("Waiting for rabbitmq to warm up...", "secondsLeft", i)
		time.Sleep(1 * time.Second)
	}
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = env.Shutdown()
	Expect(err).To(BeNil())
})
