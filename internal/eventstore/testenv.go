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

package eventstore

import (
	"context"
	"net"

	api "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/rabbitmq"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/test"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/messaging"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/storage"
	"github.com/finleap-connect/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener      net.Listener
	MetricsListener  net.Listener
	grpcServer       *grpc.Server
	messagingTestEnv *rabbitmq.TestEnv
	storageTestEnv   *storage.TestEnv
	publisher        es.EventBusPublisher
}

func (t *TestEnv) GetMessagingTestEnv() *rabbitmq.TestEnv {
	return t.messagingTestEnv
}

func NewTestEnvWithParent(testEnv *test.TestEnv) (*TestEnv, error) {
	var err error

	env := &TestEnv{
		TestEnv: testEnv,
	}

	env.messagingTestEnv, err = rabbitmq.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	conf, err := messaging.NewRabbitEventBusConfig("eventstore", env.messagingTestEnv.AmqpURL, "")
	if err != nil {
		return nil, err
	}

	env.publisher, err = messaging.NewRabbitEventBusPublisher(conf)
	if err != nil {
		return nil, err
	}

	env.storageTestEnv, err = storage.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	err = env.storageTestEnv.Store.Open(context.Background())
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServer("eventstore_grpc", false)
	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterEventStoreServer(s, NewApiServer(env.storageTestEnv.Store, env.publisher))
	})

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	env.MetricsListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// Start server
	go func() {
		err := env.grpcServer.ServeFromListener(env.apiListener, env.MetricsListener)
		if err != nil {
			panic(err)
		}
	}()

	return env, nil
}

func (env *TestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *TestEnv) Shutdown() error {
	if err := env.publisher.Close(); err != nil {
		return err
	}

	if err := env.messagingTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.storageTestEnv.Shutdown(); err != nil {
		return err
	}

	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return nil
}
