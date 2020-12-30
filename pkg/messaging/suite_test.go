package messaging

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var (
	env *MessageBusTestEnv
)

func TestMessageBus(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/messaging-junit.xml")
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

	// Start single node crdb
	container, err := env.RunWithOptions(&dockertest.RunOptions{
		Name:       "rabbitmq",
		Repository: "gitlab.figo.systems/platform/dependency_proxy/containers/bitnami/rabbitmq",
		Tag:        "3.8.9",
	})
	Expect(err).ToNot(HaveOccurred())

	// create rabbit conn
	amqpURL := fmt.Sprintf("amqp://user:bitnami@%s:%s", "127.0.0.1", container.GetPort("5672/tcp"))
	err = env.Retry(func() error {
		env.Log.Info("Connecting to rabbitmq at " + amqpURL + " ...")
		conn, err := amqp.Dial(amqpURL)
		if err != nil {
			time.Sleep(3 * time.Second)
			return err
		}
		env.RabbitConn = conn
		return nil
	})
	Expect(err).ToNot(HaveOccurred())

	publisher, err := NewRabbitEventBusPublisher(env.RabbitConn, "test")
	Expect(err).ToNot(HaveOccurred())
	env.Publisher = publisher

	consumer, err := NewRabbitEventBusConsumer(env.RabbitConn, "test")
	Expect(err).ToNot(HaveOccurred())
	env.Consumer = consumer
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())
})
