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
	})
})
