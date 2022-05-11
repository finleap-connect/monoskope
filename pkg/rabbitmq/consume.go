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
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer allows you to create and connect to queues for data consumption.
type Consumer struct {
	chManager *channelManager
	logger    logger.Logger
}

// NewConsumer returns a new Consumer connected to the given rabbitmq server
func NewConsumer(url string, config *amqp.Config) (*Consumer, error) {
	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 60 * time.Second

	var chManager *channelManager

	err := backoff.Retry(func() error {
		var localErr error

		chManager, localErr = newChannelManager(url, config, 0)

		return localErr
	}, params)
	if err != nil {
		return nil, err
	}
	consumer := Consumer{
		chManager: chManager,
		logger:    logger.WithName("rabbitmq-consumer"),
	}

	return &consumer, nil
}

// StartConsuming starts n goroutines where n="ConsumeOptions.QosOptions.Concurrency".
// Each goroutine spawns a handler that consumes off of the given queue which binds to the routing key(s).
// The provided handler is called once for each message. If the provided queue doesn't exist, it
// will be created on the cluster
func (consumer Consumer) StartConsuming(
	handler func(d amqp.Delivery) bool,
	queue string,
	routingKeys []string,
	optionFuncs ...func(*ConsumeOptions),
) error {
	defaultOptions := getDefaultConsumeOptions()
	options := &ConsumeOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.Concurrency < 1 {
		options.Concurrency = defaultOptions.Concurrency
	}

	notifyCancelOrCloseChan := make(chan error)
	consumer.chManager.registerNotify(notifyCancelOrCloseChan)

	err := consumer.startGoroutines(
		handler,
		queue,
		routingKeys,
		*options,
	)
	if err != nil {
		return err
	}

	go func() {
		for err := range notifyCancelOrCloseChan {
			consumer.logger.Info("consume cancel/close handler triggered", "error", err)
			consumer.startGoroutinesWithRetries(
				handler,
				queue,
				routingKeys,
				*options,
			)
		}
	}()
	return nil
}

// Disconnect disconnects both the channel and the connection.
// This method doesn't throw a reconnect, and should be used when finishing a program.
// IMPORTANT: If this method is executed before StopConsuming, it could cause unexpected behavior
// such as messages being processed, but not being acknowledged, thus being requeued by the broker
func (consumer Consumer) Disconnect() {
	consumer.chManager.stop()
}

// StopConsuming stops the consumption of messages.
// The consumer should be discarded as it's not safe for re-use.
// This method sends a basic.cancel notification.
// The consumerName is the name or delivery tag of the amqp consumer we want to cancel.
// When noWait is true, do not wait for the server to acknowledge the cancel.
// Only use this when you are certain there are no deliveries in flight that
// require an acknowledgment, otherwise they will arrive and be dropped in the
// client without an ack, and will not be redelivered to other consumers.
// IMPORTANT: Since the amqp library doesn't provide a way to retrieve the consumer's tag after the creation
// it's imperative for you to set the name when creating the consumer, if you want to use this function later
// a simple uuid4 should do the trick, since it should be unique.
// If you start many consumers, you should store the name of the consumers when creating them, such that you can
// use them in a for to stop all the consumers.
func (consumer Consumer) StopConsuming(consumerName string, noWait bool) {
	_ = consumer.chManager.channel.Cancel(consumerName, noWait)
}

// startGoroutinesWithRetries attempts to start consuming on a channel
// with an exponential backoff
func (consumer Consumer) startGoroutinesWithRetries(
	handler func(d amqp.Delivery) bool,
	queue string,
	routingKeys []string,
	consumeOptions ConsumeOptions,
) {
	params := backoff.NewExponentialBackOff()
	_ = backoff.Retry(func() error {
		err := consumer.startGoroutines(
			handler,
			queue,
			routingKeys,
			consumeOptions,
		)
		if err != nil {
			consumer.logger.Info("couldn't start consumer goroutines", "error", err)
		}
		return err
	}, params)
}

// startGoroutines declares the queue if it doesn't exist,
// binds the queue to the routing key(s), and starts the goroutines
// that will consume from the queue
func (consumer Consumer) startGoroutines(
	handler func(d amqp.Delivery) bool,
	queue string,
	routingKeys []string,
	consumeOptions ConsumeOptions,
) error {
	consumer.chManager.channelMux.RLock()
	defer consumer.chManager.channelMux.RUnlock()

	_, err := consumer.chManager.channel.QueueDeclare(
		queue,
		consumeOptions.QueueDurable,
		consumeOptions.QueueAutoDelete,
		consumeOptions.QueueExclusive,
		consumeOptions.QueueNoWait,
		consumeOptions.QueueArgs,
	)
	if err != nil {
		return err
	}

	if consumeOptions.BindingExchange != nil {
		exchange := consumeOptions.BindingExchange
		if exchange.Name == "" {
			return fmt.Errorf("binding to exchange but name not specified")
		}
		err = consumer.chManager.channel.ExchangeDeclare(
			exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.ExchangeArgs,
		)
		if err != nil {
			return err
		}
		for _, routingKey := range routingKeys {
			err = consumer.chManager.channel.QueueBind(
				queue,
				routingKey,
				exchange.Name,
				consumeOptions.BindingNoWait,
				consumeOptions.BindingArgs,
			)
			if err != nil {
				return err
			}
		}
	}

	err = consumer.chManager.channel.Qos(
		consumeOptions.QOSPrefetch,
		0,
		consumeOptions.QOSGlobal,
	)
	if err != nil {
		return err
	}

	msgs, err := consumer.chManager.channel.Consume(
		queue,
		consumeOptions.ConsumerName,
		consumeOptions.ConsumerAutoAck,
		consumeOptions.ConsumerExclusive,
		consumeOptions.ConsumerNoLocal, // no-local is not supported by RabbitMQ
		consumeOptions.ConsumerNoWait,
		consumeOptions.ConsumerArgs,
	)
	if err != nil {
		return err
	}

	for i := 0; i < consumeOptions.Concurrency; i++ {
		go func() {
			for msg := range msgs {
				if consumeOptions.ConsumerAutoAck {
					handler(msg)
					continue
				}
				if handler(msg) {
					err := msg.Ack(false)
					if err != nil {
						consumer.logger.Error(err, "can't ack message")
					}
				} else {
					err := msg.Nack(false, true)
					if err != nil {
						consumer.logger.Error(err, "can't nack message")
					}
				}
			}
			consumer.logger.Info("rabbit consumer goroutine closed")
		}()
	}
	consumer.logger.Info("Processing messages on goroutines", "concurrency", consumeOptions.Concurrency)
	return nil
}
