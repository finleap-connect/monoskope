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

package storage

import (
	"errors"
	"strings"
	"time"

	m8tls "github.com/finleap-connect/monoskope/pkg/tls"

	"github.com/go-pg/pg/v10"
)

const (
	DefaultReconnectDelay = 10 * time.Second // When reconnecting to the server after connection failure
	DefaultReInitDelay    = 5 * time.Second  // When setting up db schema
	DefaultResendDelay    = 3 * time.Second  // When retrying to read/write
	DefaultMaxRetries     = 10               // How many times retrying read/write
	CACertPath            = "/etc/eventstore/certs/db/ca.crt"
	TLSCertPath           = "/etc/eventstore/certs/db/tls.crt"
	TLSKeyPath            = "/etc/eventstore/certs/db/tls.key"
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
	loader, err := m8tls.NewTLSConfigLoader()
	if err != nil {
		return err
	}

	err = loader.SetServerCACertificate(CACertPath)
	if err != nil {
		return err
	}

	err = loader.SetClientCertificate(TLSCertPath, TLSKeyPath)
	if err != nil {
		return err
	}

	err = loader.Watch()
	if err != nil {
		return err
	}

	conf.pgOptions.TLSConfig = loader.GetClientTLSConfig()
	conf.pgOptions.TLSConfig.ServerName = strings.Split(conf.pgOptions.Addr, ":")[0]

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
