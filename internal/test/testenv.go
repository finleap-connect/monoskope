package test

import (
	"fmt"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
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

func (t *TestEnv) RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*dc.HostConfig)) (*dockertest.Resource, error) {
	t.Log.Info(fmt.Sprintf("Starting docker container %s:%s ...", opts.Repository, opts.Tag))
	res, err := t.pool.RunWithOptions(opts, hcOpts...)
	if err != nil {
		return nil, err
	}
	t.resources[res.Container.Name] = res

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
