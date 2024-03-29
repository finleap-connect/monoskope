// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rabbitmq

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ = Describe("Pkg/Rabbitmq/Consume", func() {
	expectedExchangeName := "test-exchange"
	expectedRoutingKey := "*"
	expectedConsumerName := "test-consumer"

	It("can consume", func() {
		consumer, err := NewConsumer(env.AmqpURL, &amqp.Config{})
		Expect(err).ToNot(HaveOccurred())
		Expect(consumer).ToNot(BeNil())

		publisher, _, err := NewPublisher(env.AmqpURL, &amqp.Config{})
		Expect(err).ToNot(HaveOccurred())
		defer publisher.StopPublishing()

		done := make(chan interface{})

		err = consumer.StartConsuming(func(d amqp.Delivery) bool {
			defer close(done)
			Expect(d).ToNot(BeNil())
			return true
		}, "test-consumer-queue", []string{expectedRoutingKey},
			WithConsumeOptionsConsumerName(expectedConsumerName),
			WithConsumeOptionsBindingExchangeName(expectedExchangeName),
			WithConsumeOptionsBindingExchangeKind(amqp.ExchangeTopic),
			WithConsumeOptionsBindingExchangeDurable,
		)

		Expect(err).ToNot(HaveOccurred())
		defer consumer.StopConsuming(expectedConsumerName, false)

		err = publisher.Publish(context.Background(), []byte("test"), []string{expectedRoutingKey}, WithPublishOptionsExchange(expectedExchangeName))
		Expect(err).ToNot(HaveOccurred())

		Eventually(done, 60).Should(BeClosed())
	})
})
