package rabbitmq

import (
	"errors"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

type channelManager struct {
	logger              logger.Logger
	url                 string
	channel             *amqp.Channel
	connection          *amqp.Connection
	config              *amqp.Config
	channelMux          *sync.RWMutex
	notifyCancelOrClose chan error
	maximumBackoff      time.Duration
}

func newChannelManager(url string, conf *amqp.Config, maximumBackoff time.Duration) (*channelManager, error) {
	conn, ch, err := getNewChannel(url, conf)
	if err != nil {
		return nil, err
	}

	chManager := channelManager{
		logger:              logger.WithName("rabbitmq-channel-manager"),
		url:                 url,
		connection:          conn,
		channel:             ch,
		channelMux:          &sync.RWMutex{},
		config:              conf,
		notifyCancelOrClose: make(chan error),
		maximumBackoff:      maximumBackoff,
	}
	go chManager.startNotifyCancelOrClosed()
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
	notifyCloseConn := chManager.connection.NotifyClose(make(chan *amqp.Error, 1))

	select {
	case err := <-notifyCloseChan:
		// If the connection close is triggered by the Server, a reconnection takes place
		if err != nil && err.Server {
			chManager.logger.Info("attempting to reconnect to amqp server after channel close")
			chManager.reconnectWithBackoff()
			chManager.logger.Info("successfully reconnected to amqp server after channel close")
			chManager.notifyCancelOrClose <- err
		}
	case err := <-notifyCloseConn:
		// If the connection close is triggered by the Server, a reconnection takes place
		if err != nil && err.Server {
			chManager.logger.Info("attempting to reconnect to amqp server after connection close")
			chManager.reconnectWithBackoff()
			chManager.logger.Info("successfully reconnected to amqp server after connection close")
			chManager.notifyCancelOrClose <- err
		}
	case err := <-notifyCancelChan:
		chManager.logger.Info("attempting to reconnect to amqp server after cancel")
		chManager.reconnectWithBackoff()
		chManager.logger.Info("successfully reconnected to amqp server after cancel")
		chManager.notifyCancelOrClose <- errors.New(err)
	}
}

// reconnectWithBackoff continuously attempts to reconnect with an
// exponential backoff strategy
func (chManager *channelManager) reconnectWithBackoff() {
	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = chManager.maximumBackoff
	err := backoff.Retry(func() error {
		reconnectErr := chManager.reconnect()
		if reconnectErr != nil {
			chManager.logger.Info("waiting to attempt to reconnect to amqp server", "error", reconnectErr.Error())
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
