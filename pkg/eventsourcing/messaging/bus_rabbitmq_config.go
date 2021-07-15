package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

const (
	DefaultExchangeName = "m8_events" // name of the monoskope exchange
	CACertPath          = "/etc/eventstore/certs/buscerts/ca.crt"
	TLSCertPath         = "/etc/eventstore/certs/buscerts/tls.crt"
	TLSKeyPath          = "/etc/eventstore/certs/buscerts/tls.key"
)

type RabbitEventBusConfig struct {
	name             string // Name of the client, required
	url              string // Connection string, required
	routingKeyPrefix string // Prefix for routing of messages
	exchangeName     string // Name of the exchange to initialize/use
	amqpConfig       *amqp.Config
}

// NewRabbitEventBusConfig creates a new RabbitEventBusConfig with defaults.
func NewRabbitEventBusConfig(name, url, routingKeyPrefix string) *RabbitEventBusConfig {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}

	return &RabbitEventBusConfig{
		name:             name,
		url:              url,
		routingKeyPrefix: routingKeyPrefix,
		exchangeName:     DefaultExchangeName,
		amqpConfig:       &amqp.Config{},
	}
}

// ConfigureTLS adds the configuration for TLS secured connection/auth
func (conf *RabbitEventBusConfig) ConfigureTLS() error {
	conf.amqpConfig.SASL = []amqp.Authentication{&CertAuth{}}
	if err := loadCertificates(conf.amqpConfig); err != nil {
		return err
	}
	return nil
}

func loadCertificates(amqpConfig *amqp.Config) error {
	cfg := &tls.Config{
		RootCAs: x509.NewCertPool(),
	}

	if ca, err := ioutil.ReadFile(CACertPath); err != nil {
		return err
	} else {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	if cert, err := tls.LoadX509KeyPair(TLSCertPath, TLSKeyPath); err != nil {
		return err
	} else {
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	amqpConfig.TLSClientConfig = cfg

	return nil
}

// Validate validates the configuration
func (conf *RabbitEventBusConfig) Validate() error {
	if conf.name == "" {
		return errors.ErrConfigNameRequired
	}
	if conf.url == "" {
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
