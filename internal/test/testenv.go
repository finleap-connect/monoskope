package test

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

type TestEnv struct {
	pool     *dockertest.Pool
	Log      logger.Logger
	shutdown bool
}

func IsRunningInCI() bool {
	_, runningInCi := os.LookupEnv("CI")
	return runningInCi
}

func (t *TestEnv) CreateDockerPool() error {
	// Running in CI, no docker necessary
	if IsRunningInCI() {
		return nil
	}

	t.Log.Info("Creating docker pool...")

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	t.pool = pool
	return nil
}

func (t *TestEnv) Retry(op func() error) error {
	return t.pool.Retry(op)
}

func (t *TestEnv) Run(opts *dockertest.RunOptions) (*dockertest.Resource, error) {
	res, present := t.pool.ContainerByName(opts.Name)
	if present {
		return res, nil
	}

	t.Log.Info(fmt.Sprintf("Starting docker container %s:%s ...", opts.Repository, opts.Tag))
	res, err := t.pool.RunWithOptions(opts, func(config *dc.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = dc.NeverRestart()
	})
	if err != nil {
		return nil, err
	}

	err = res.Expire(60)
	if err != nil {
		return nil, err
	}

	containerLogger := logWriter{}
	logOptions := dc.LogsOptions{
		Container:    opts.Name,
		Follow:       true,
		OutputStream: containerLogger,
		ErrorStream:  containerLogger,
		Stdout:       true,
		Stderr:       true,
		Context:      context.Background(),
	}
	go func() {
		err = t.pool.Client.Logs(logOptions)
		if err != nil {
			t.Log.Error(err, err.Error())
		}
	}()

	return res, err
}

func NewTestEnv(envName string) *TestEnv {
	log := logger.WithName(envName)
	env := &TestEnv{
		Log: log,
	}
	log.Info("Setting up testenv...")
	return env
}

func (env *TestEnv) Shutdown() error {
	if env.shutdown {
		return nil
	}

	env.shutdown = true
	log := env.Log
	log.Info("Tearing down testenv...")

	return nil
}

type logWriter struct {
}

func (l logWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}
