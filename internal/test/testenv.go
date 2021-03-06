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

package test

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

type TestEnv struct {
	pool         *dockertest.Pool
	Log          logger.Logger
	shutdown     bool
	keepExisting bool
	resources    map[string]*dockertest.Resource
}

func (t *TestEnv) CreateDockerPool(keepExisting bool) error {
	t.Log.Info("Creating docker pool...")

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	t.pool = pool
	t.keepExisting = keepExisting

	return nil
}

func (t *TestEnv) Retry(op func() error) error {
	return t.pool.Retry(op)
}

func (t *TestEnv) Run(opts *dockertest.RunOptions) (*dockertest.Resource, error) {
	res, present := t.pool.ContainerByName(opts.Name)
	if present {
		if t.keepExisting {
			return res, nil
		} else {
			if err := t.pool.Purge(res); err != nil {
				return nil, err
			}
		}
	}

	t.Log.Info(fmt.Sprintf("Starting docker container %s:%s ...", opts.Repository, opts.Tag))
	res, err := t.pool.RunWithOptions(opts, func(config *dc.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = dc.NeverRestart()
	})
	if err != nil {
		return nil, err
	}
	t.resources[res.Container.Name] = res

	if t.keepExisting {
		err = res.Expire(500)
		if err != nil {
			return nil, err
		}
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

func (t *TestEnv) Purge(resource string) error {
	res, present := t.pool.ContainerByName(resource)
	if present {
		return t.pool.Purge(res)
	}
	delete(t.resources, resource)

	return nil
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
	if env.shutdown {
		return nil
	}

	if !env.keepExisting {
		if env.resources != nil {
			for key, element := range env.resources {
				env.Log.Info("Tearing down docker resource", "resource", key)
				if err := env.pool.Purge(element); err != nil {
					return err
				}
			}
		}
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
