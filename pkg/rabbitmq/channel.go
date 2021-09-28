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
	"errors"
	"io"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type channelManager struct {
	logger                       logger.Logger
	url                          string
	channel                      *amqp.Channel
	connection                   *amqp.Connection
	config                       *amqp.Config
	channelMux                   *sync.RWMutex
	notifyCancelOrCloseBroadcast []chan<- error
	isReconnecting               bool
	maximumBackoff               time.Duration
}

func newChannelManager(url string, conf *amqp.Config, maximumBackoff time.Duration) (*channelManager, error) {
	log := logger.WithName("rabbitmq-channel-manager")
	log.Info("attempting to connect to amqp server", "server", url)

	conn, ch, err := getNewChannel(url, conf)
	if err != nil {
		return nil, err
	}

	chManager := channelManager{
		logger:         log,
		url:            url,
		connection:     conn,
		channel:        ch,
		channelMux:     &sync.RWMutex{},
		config:         conf,
		maximumBackoff: maximumBackoff,
	}
	go chManager.startNotifyCancelOrClosed()
	log.Info("connected to amqp server!", "url", url)

	return &chManager, nil
}

func getNewChannel(url string, conf *amqp.Config) (*amqp.Connection, *amqp.Channel, error) {
	amqpConn, err := amqp.DialConfig(url, *conf)
	if err != nil {
		return nil, nil, err
	}
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return amqpConn, ch, err
}

// startNotifyCancelOrClosed listens on the channel's cancelled and closed
// notifiers. When it detects a problem, it attempts to reconnect with an exponential
// backoff. Once reconnected, it sends an error back on the manager's notifyCancelOrClose
// channel
func (chManager *channelManager) startNotifyCancelOrClosed() {
	notifyCloseChan := chManager.channel.NotifyClose(make(chan *amqp.Error, 1))
	notifyCancelChan := chManager.channel.NotifyCancel(make(chan string, 1))

	select {
	case err := <-notifyCloseChan:
		// If the connection close is triggered by the Server, a reconnection takes place
		if err != nil && (err.Server || err.Reason == io.EOF.Error()) {
			chManager.logger.Info("attempting to reconnect to amqp server after channel close")
			chManager.reconnectWithBackoff()
			chManager.logger.Info("successfully reconnected to amqp server after channel close")
			chManager.notifyListener(err)
		}
	case err := <-notifyCancelChan:
		chManager.logger.Info("attempting to reconnect to amqp server after cancel")
		chManager.reconnectWithBackoff()
		chManager.logger.Info("successfully reconnected to amqp server after cancel")
		chManager.notifyListener(errors.New(err))
	}
}

// reconnectWithBackoff continuously attempts to reconnect with an
// exponential backoff strategy, it never stops if maximumBackoff is set to zero.
func (chManager *channelManager) reconnectWithBackoff() {
	chManager.isReconnecting = true
	defer func() { chManager.isReconnecting = false }()

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = chManager.maximumBackoff
	err := backoff.Retry(func() error {
		reconnectErr := chManager.reconnect()
		if reconnectErr != nil {
			chManager.logger.Info("waiting to attempt to reconnect to amqp server", "backoff", params.NextBackOff(), "error", reconnectErr.Error())
		}
		return reconnectErr
	}, params)
	if err != nil {
		chManager.logger.Error(err, "backoff reached maximum backoff", "backoff", chManager.maximumBackoff)
		util.PanicOnError(err)
	}
}

// reconnect safely closes the current channel and obtains a new one
func (chManager *channelManager) reconnect() error {
	chManager.channelMux.Lock()
	defer chManager.channelMux.Unlock()
	newConn, newChannel, err := getNewChannel(chManager.url, chManager.config)
	if err != nil {
		return err
	}

	chManager.channel.Close()
	chManager.connection.Close()

	chManager.connection = newConn
	chManager.channel = newChannel
	go chManager.startNotifyCancelOrClosed()
	return nil
}

// stop closes the channel and connection
func (chManager *channelManager) stop() {
	chManager.channel.Close()
	chManager.connection.Close()
}

func (chManager *channelManager) registerNotify(errChan chan<- error) {
	chManager.notifyCancelOrCloseBroadcast = append(chManager.notifyCancelOrCloseBroadcast, errChan)
}

func (chManager *channelManager) notifyListener(err error) {
	for _, listener := range chManager.notifyCancelOrCloseBroadcast {
		listener <- err
	}
}
