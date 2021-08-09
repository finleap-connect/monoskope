package rabbitmq

import (
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/rabbitmq/amqp091-go"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv
	AmqpURL          string
	ExternalRabbitMQ bool
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

	if err := env.startRabbitMQ(); err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) stopRabbitMQ() error {
	if !env.ExternalRabbitMQ {
		env.Log.Info("Purging rabbitmq...")
		return env.Purge("rabbitmq")
	}
	return nil
}

func (env *TestEnv) startRabbitMQ() error {
	env.Log.Info("Starting rabbitmq...")
	var err error

	if v := os.Getenv("AMQP_URL"); v != "" {
		env.AmqpURL = v // running in ci pipeline
		env.ExternalRabbitMQ = true
	} else {
		// Start rabbitmq
		_, err := env.Run(&dockertest.RunOptions{
			Name:       "rabbitmq",
			Repository: "artifactory.figo.systems/public_docker/bitnami/rabbitmq",
			Tag:        "3.8.19",
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
		env.AmqpURL = "amqp://user:bitnami@127.0.0.1:5672"
	}

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
