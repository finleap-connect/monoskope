// Copyright 2021 Monoskope Authors
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
