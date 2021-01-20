package test

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	_ "go.uber.org/automaxprocs"
)

type TestEnv struct {
	pool      *dockertest.Pool
	resources map[string]*dockertest.Resource
	Log       logger.Logger
}

func (t *TestEnv) CreateDockerPool() error {
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
		if err := t.pool.Purge(res); err != nil {
			return nil, err
		}
	}

	t.Log.Info(fmt.Sprintf("Starting docker container %s:%s ...", opts.Repository, opts.Tag))
	res, err := t.pool.RunWithOptions(opts)
	if err != nil {
		return nil, err
	}
	t.resources[res.Container.Name] = res

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
		Log:       log,
		resources: make(map[string]*dockertest.Resource),
	}
	log.Info("Setting up testenv...")
	return env
}

func (env *TestEnv) Shutdown() error {
	log := env.Log
	log.Info("Tearing down testenv...")

	if env.resources != nil {
		for key, element := range env.resources {
			log.Info("Tearing down docker resource", "resource", key)
			if err := env.pool.Purge(element); err != nil {
				return err
			}
		}
	}

	return nil
}

type logWriter struct {
}

func (l logWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}
