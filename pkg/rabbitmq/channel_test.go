package rabbitmq

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rabbitmq/amqp091-go"
)

var _ = Describe("Pkg/Rabbitmq/Channel", func() {
	It("can connect", func() {
		chanManager, err := newChannelManager(env.AmqpURL, &amqp091.Config{}, time.Second*5)
		Expect(err).ToNot(HaveOccurred())
		Expect(chanManager).ToNot(BeNil())
		Expect(chanManager.channel.IsClosed()).To(BeFalse())
	})
})
