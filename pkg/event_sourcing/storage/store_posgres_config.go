package storage

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

const (
	DefaultReconnectDelay = 10 * time.Second // When reconnecting to the server after connection failure
	DefaultReInitDelay    = 5 * time.Second  // When setting up db schema
	DefaultResendDelay    = 3 * time.Second  // When retrying to read/write
	DefaultMaxRetries     = 10               // How many times retrying read/write
)

type postgresStoreConfig struct {
	ReconnectDelay time.Duration // When reconnecting to the server after connection failure
	ReInitDelay    time.Duration // When setting up db schema
	RetryDelay     time.Duration // When retrying to read/write
	MaxRetries     int           // How many times retrying read/write
	pgOptions      *pg.Options
}

// ErrConfigDbNameRequired is when the config doesn't include a name.
var ErrConfigDbNameRequired = errors.New("database name must not be empty")

// ErrConfigUrlRequired is when the config doesn't include a name.
var ErrConfigUrlRequired = errors.New("url must not be empty")

// NewRabbitEventBusConfig creates a new RabbitEventBusConfig with defaults.
func NewPostgresStoreConfig(url string) (*postgresStoreConfig, error) {
	options, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return &postgresStoreConfig{
		ReconnectDelay: DefaultReconnectDelay,
		RetryDelay:     DefaultResendDelay,
		MaxRetries:     DefaultMaxRetries,
		ReInitDelay:    DefaultReInitDelay,
		pgOptions:      options,
	}, nil
}

// ConfigureTLS adds the configuration for TLS secured connection/auth
func (conf *postgresStoreConfig) ConfigureTLS() error {
	cfg := &tls.Config{
		RootCAs:    x509.NewCertPool(),
		ServerName: strings.Split(conf.pgOptions.Addr, ":")[0],
	}
	if ca, err := ioutil.ReadFile("/etc/eventstore/certs/db/ca.crt"); err != nil {
		return err
	} else {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	if cert, err := tls.LoadX509KeyPair("/etc/eventstore/certs/db/tls.crt", "/etc/eventstore/certs/db/tls.key"); err != nil {
		return err
	} else {
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	conf.pgOptions.TLSConfig = cfg
	return nil
}

// Validate validates the configuration
func (conf *postgresStoreConfig) Validate() error {
	if conf.pgOptions.Database == "" {
		return ErrConfigDbNameRequired
	}
	if conf.pgOptions.Addr == "" {
		return ErrConfigUrlRequired
	}
	return nil
}
