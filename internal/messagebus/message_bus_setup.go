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

package messagebus

import (
	"os"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esMessaging "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
)

func getMsgBusUrl() string {
	var msgbusUrl string
	if v := os.Getenv("BUS_URL"); v != "" {
		msgbusUrl = v
	}
	return msgbusUrl
}

func getRabbitConf(name, msgbusPrefix string) (*esMessaging.RabbitEventBusConfig, error) {
	rabbitConf, err := esMessaging.NewRabbitEventBusConfig(name, getMsgBusUrl(), msgbusPrefix)
	if err != nil {
		return nil, err
	}
	return rabbitConf, nil
}

func NewEventBusConsumer(name, msgbusPrefix string) (eventsourcing.EventBusConsumer, error) {
	rabbitConf, err := getRabbitConf(name, msgbusPrefix)
	if err != nil {
		return nil, err
	}
	return NewEventBusConsumerFromConfig(rabbitConf)
}

func NewEventBusConsumerFromConfig(rabbitConf *esMessaging.RabbitEventBusConfig) (eventsourcing.EventBusConsumer, error) {
	consumer, err := esMessaging.NewRabbitEventBusConsumer(rabbitConf)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func NewEventBusPublisher(name, msgbusPrefix string) (eventsourcing.EventBusPublisher, error) {
	rabbitConf, err := getRabbitConf(name, msgbusPrefix)
	if err != nil {
		return nil, err
	}

	publisher, err := esMessaging.NewRabbitEventBusPublisher(rabbitConf)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
