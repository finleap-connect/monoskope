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

package messaging

import (
	m8tls "github.com/finleap-connect/monoskope/pkg/tls"

	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultExchangeName = "m8_events" // name of the monoskope exchange
	CACertPath          = "/etc/eventstore/certs/buscerts/ca.crt"
	TLSCertPath         = "/etc/eventstore/certs/buscerts/tls.crt"
	TLSKeyPath          = "/etc/eventstore/certs/buscerts/tls.key"
)

type RabbitEventBusConfig struct {
	url              string // Connection string, required
	Name             string // Name of the client, required
	RoutingKeyPrefix string // Prefix for routing of messages
	ExchangeName     string // Name of the exchange to initialize/use
	AMQPConfig       *amqp.Config
}

// NewRabbitEventBusConfig creates a new RabbitEventBusConfig with defaults.
func NewRabbitEventBusConfig(name, url, routingKeyPrefix string) (*RabbitEventBusConfig, error) {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}

	conf := &RabbitEventBusConfig{
		Name:             name,
		url:              url,
		RoutingKeyPrefix: routingKeyPrefix,
		ExchangeName:     DefaultExchangeName,
		AMQPConfig:       &amqp.Config{},
	}

	if err := conf.SetURL(url); err != nil {
		return nil, err
	}
	return conf, nil
}

// URL of the RabbitMQ host to connect to
func (conf *RabbitEventBusConfig) URL() string {
	return conf.url
}

// SetURL reconfigures the address of the RabbitMQ host
func (conf *RabbitEventBusConfig) SetURL(url string) error {
	uri, err := amqp.ParseURI(url)
	if err != nil {
		return err
	}
	if uri.Scheme == "amqps" {
		if err := conf.configureTLS(); err != nil {
			return err
		}
	}
	conf.url = url
	return nil
}

// configureTLS adds the configuration for TLS secured connection/auth
func (conf *RabbitEventBusConfig) configureTLS() error {
	loader, err := m8tls.NewTLSConfigLoader(CACertPath, TLSCertPath, TLSKeyPath)
	if err != nil {
		return err
	}

	err = loader.Watch()
	if err != nil {
		return err
	}

	conf.AMQPConfig.TLSClientConfig = loader.GetTLSConfig()
	conf.AMQPConfig.SASL = []amqp.Authentication{&CertAuth{}}

	return nil
}

// Validate validates the configuration
func (conf *RabbitEventBusConfig) Validate() error {
	if conf.Name == "" {
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
