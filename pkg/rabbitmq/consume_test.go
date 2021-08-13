package rabbitmq

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rabbitmq/amqp091-go"
)

var _ = Describe("Pkg/Rabbitmq/Consume", func() {
	It("can consume", func() {
		consumer, err := NewConsumer(env.AmqpURL, &amqp091.Config{})
		Expect(err).ToNot(HaveOccurred())
		Expect(consumer).ToNot(BeNil())

		err = consumer.StartConsuming(func(d amqp091.Delivery) bool {
			Expect(d).ToNot(BeNil())
			return true
		}, "test-consumer-queue", []string{"*"}, WithConsumeOptionsConsumerName("test-consumer"))
		Expect(err).ToNot(HaveOccurred())

		publisher, _, err := NewPublisher(env.AmqpURL, &amqp091.Config{})
		Expect(err).ToNot(HaveOccurred())

		err = publisher.Publish([]byte("test"), []string{"*"})
		Expect(err).ToNot(HaveOccurred())

		consumer.StopConsuming("test-consumer", false)
	})
})
