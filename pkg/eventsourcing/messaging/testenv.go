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

	warumupSeconds := 30
	if test.IsRunningInCI() {
		warumupSeconds = 0 // no warmup necessary in CI
	}

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

		// Build connection string
		env.AmqpURL = fmt.Sprintf("amqp://user:bitnami@127.0.0.1:%s", container.GetPort("5672/tcp"))
	}

	// Wait for rabbitmq to start
	for i := warumupSeconds; i > 0; i-- {
		env.Log.Info("Waiting for rabbitmq to warm up...", "secondsLeft", i)
		time.Sleep(1 * time.Second)
	}

	return env, nil
}
