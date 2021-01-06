package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"github.com/streadway/amqp"
)

const (
	DefaultExchangeName   = "m8_events"      // name of the monoskope exchange
	DefaultReconnectDelay = 10 * time.Second // When reconnecting to the server after connection failure
	DefaultReInitDelay    = 5 * time.Second  // When setting up the channel after a channel exception
	DefaultResendDelay    = 3 * time.Second  // When resending messages the server didn't confirm
	DefaultMaxResends     = 10               // How many times resending messages the server didn't confirm
)

type rabbitEventBusConfig struct {
	Name             string        // Name of the client, required
	Url              string        // Connection string, required
	RoutingKeyPrefix string        // Prefix for routing of messages
	ExchangeName     string        // Name of the exchange to initialize/use
	ReconnectDelay   time.Duration // When reconnecting to the server after connection failure
	ResendDelay      time.Duration // When resending messages the server didn't confirm
	MaxResends       int           // How many times resending messages the server didn't confirm
	ReInitDelay      time.Duration // When setting up the channel after a channel exception
	amqpConfig       amqp.Config
}

// ErrConfigNameRequired is when the config doesn't include a name.
var ErrConfigNameRequired = errors.New("name must not be empty")

// ErrConfigUrlRequired is when the config doesn't include a name.
var ErrConfigUrlRequired = errors.New("url must not be empty")

// NewRabbitEventBusConfig creates a new RabbitEventBusConfig with defaults.
func NewRabbitEventBusConfig(name, url string) *rabbitEventBusConfig {
	return &rabbitEventBusConfig{
		Name:             name,
		Url:              url,
		RoutingKeyPrefix: "m8",
		ExchangeName:     DefaultExchangeName,
		ReconnectDelay:   DefaultReconnectDelay,
		ResendDelay:      DefaultResendDelay,
		MaxResends:       DefaultMaxResends,
		ReInitDelay:      DefaultReInitDelay,
		amqpConfig:       amqp.Config{},
	}
}

// ConfigureTLS adds the configuration for TLS secured connection/auth
func (conf *rabbitEventBusConfig) ConfigureTLS() error {
	cfg := &tls.Config{
		RootCAs: x509.NewCertPool(),
	}
	if ca, err := ioutil.ReadFile("/etc/eventstore/certs/bus/ca.crt"); err != nil {
		return err
	} else {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	if cert, err := tls.LoadX509KeyPair("/etc/eventstore/certs/bus/tls.crt", "/etc/eventstore/certs/bus/tls.key"); err != nil {
		return err
	} else {
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	conf.amqpConfig.TLSClientConfig = cfg
	return nil
}

// Validate validates the configuration
func (conf *rabbitEventBusConfig) Validate() error {
	if conf.Name == "" {
		return ErrConfigNameRequired
	}
	if conf.Url == "" {
		return ErrConfigUrlRequired
	}
	return nil
}
