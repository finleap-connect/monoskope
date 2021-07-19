package rabbitmq

import (
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/ory/dockertest/v3"
	"github.com/rabbitmq/amqp091-go"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
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

	if err := env.CreateDockerPool(true); err != nil {
		return nil, err
	}

	if v := os.Getenv("AMQP_URL"); v != "" {
		env.AmqpURL = v // running in ci pipeline
	} else {
		// Start rabbitmq
		container, err := env.Run(&dockertest.RunOptions{
			Name:       "rabbitmq",
			Repository: "artifactory.figo.systems/public_docker/bitnami/rabbitmq",
			Tag:        "3.8.19",
			Env: []string{
				"RABBITMQ_PLUGINS=rabbitmq_management",
			},
		})
		if err != nil {
			return nil, err
		}

		// Build connection string
		env.AmqpURL = fmt.Sprintf("amqp://user:bitnami@127.0.0.1:%s", container.GetPort("5672/tcp"))
	}

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 60 * time.Second
	err := backoff.Retry(func() error {
		cm, err := newChannelManager(env.AmqpURL, &amqp091.Config{}, time.Second*5)
		if err == nil {
			cm.stop()
		}
		return err
	}, params)

	return env, err
}
