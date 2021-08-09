package rabbitmq

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rabbitmq/amqp091-go"
)

var _ = Describe("Pkg/Rabbitmq/Channel", func() {
	It("can connect with reconnect", func() {
		chanManager, err := newChannelManager(env.AmqpURL, &amqp091.Config{}, 0)
		Expect(err).ToNot(HaveOccurred())
		Expect(chanManager).ToNot(BeNil())
		Expect(chanManager.channel.IsClosed()).To(BeFalse())

		if !env.ExternalRabbitMQ {
			err = env.stopRabbitMQ()
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(3 * time.Second)

			Expect(chanManager.channel.IsClosed()).To(BeTrue())
			Expect(chanManager.isReconnecting).To(BeTrue())

			err = env.startRabbitMQ()
			Expect(err).ToNot(HaveOccurred())

			for chanManager.isReconnecting {
				time.Sleep(300 * time.Millisecond)
			}
			Expect(chanManager.channel.IsClosed()).To(BeFalse())
		}
	})
})
