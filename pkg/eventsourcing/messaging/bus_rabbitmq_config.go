package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

const (
	DefaultExchangeName   = "m8_events" // name of the monoskope exchange
	DefaultHeartbeat      = 10 * time.Second
	DefaultReconnectDelay = 10 * time.Second // When reconnecting to the server after connection failure
	DefaultReInitDelay    = 5 * time.Second  // When setting up the channel after a channel exception
	DefaultResendDelay    = 5 * time.Second  // When resending messages the server didn't confirm
	DefaultMaxResends     = 5                // How many times resending messages the server didn't confirm
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
	amqpConfig       *amqp.Config
}

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
		amqpConfig:       &amqp.Config{},
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

	conf.amqpConfig.Heartbeat = DefaultHeartbeat
	conf.amqpConfig.TLSClientConfig = cfg
	conf.amqpConfig.SASL = []amqp.Authentication{&CertAuth{}}

	return nil
}

// Validate validates the configuration
func (conf *rabbitEventBusConfig) Validate() error {
	if conf.Name == "" {
		return errors.ErrConfigNameRequired
	}
	if conf.Url == "" {
		return errors.ErrConfigUrlRequired
	}
	return nil
}

// CertAuth for RabbitMQ-auth-mechanism-ssl.
type CertAuth struct {
}

func (me *CertAuth) Mechanism() string {
	return "EXTERNAL"
}

func (me *CertAuth) Response() string {
	return "\000*\000*"
}
