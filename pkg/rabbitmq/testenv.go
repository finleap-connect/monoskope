// Copyright 2022 Monoskope Authors
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
	"time"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/rabbitmq/amqp091-go"
)

type TestEnv struct {
	*test.TestEnv
	AmqpURL string
}

func (env *TestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func NewTestEnvWithParent(testEnv *test.TestEnv) (*TestEnv, error) {
	env := &TestEnv{
		TestEnv: testEnv,
	}

	if err := env.CreateDockerPool(false); err != nil {
		return nil, err
	}

	if err := env.startRabbitMQ(); err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) stopRabbitMQ() error {
	env.Log.Info("Purging rabbitmq...")
	return env.Purge("rabbitmq")
}

func (env *TestEnv) startRabbitMQ() error {
	env.Log.Info("Starting rabbitmq...")
	var err error

	// Start rabbitmq
	_, err = env.Run(&dockertest.RunOptions{
		Name:       "rabbitmq",
		Repository: "rabbitmq",
		Tag:        "3.10.2",
		Env: []string{
			"RABBITMQ_PLUGINS=rabbitmq_management",
		},
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5672/tcp": {{HostPort: "5672"}},
		},
	})
	if err != nil {
		return err
	}

	// Build connection string
	env.AmqpURL = "amqp://guest:guest@127.0.0.1:5672"

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 60 * time.Second
	err = backoff.Retry(func() error {
		cm, err := newChannelManager(env.AmqpURL, &amqp091.Config{}, 0)
		if err == nil {
			cm.stop()
		}
		return err
	}, params)

	if err != nil {
		env.Log.Error(err, "Starting rabbitmq failed!")
	} else {
		env.Log.Info("Started rabbitmq!")
	}
	return err
}
