package messaging

import (
	"fmt"
	"os"
	"time"

	"github.com/ory/dockertest/v3"
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

	if err := env.CreateDockerPool(); err != nil {
		return nil, err
	}

	rabbitWarmUpSeconds := 30

	if v := os.Getenv("AMQP_URL"); v != "" {
		env.AmqpURL = v // running in ci pipeline
	} else {
		// Start rabbitmq
		container, err := env.Run(&dockertest.RunOptions{
			Name:       "rabbitmq",
			Repository: "artifactory.figo.systems/public_docker/bitnami/rabbitmq",
			Tag:        "3.8.9",
			Env: []string{
				"RABBITMQ_PLUGINS=rabbitmq_management",
			},
		})
		if err != nil {
			return nil, err
		}
		if container.Container.Created.Before(time.Now().Add(time.Second * time.Duration(-rabbitWarmUpSeconds))) {
			rabbitWarmUpSeconds = 0
		}

		// Build connection string
		env.AmqpURL = fmt.Sprintf("amqp://user:bitnami@127.0.0.1:%s", container.GetPort("5672/tcp"))
	}

	// Wait for rabbitmq to start
	if test.IsRunningInCI() {
		rabbitWarmUpSeconds = 0 // no warmup necessary in CI
	}

	for i := rabbitWarmUpSeconds; i > 0; i-- {
		env.Log.Info("Waiting for rabbitmq to warm up...", "secondsLeft", i)
		time.Sleep(1 * time.Second)
	}

	return env, nil
}
