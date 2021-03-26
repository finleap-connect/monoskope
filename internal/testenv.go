package internal

import (
	"os"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv
	eventStoreTestEnv     *eventstore.TestEnv
	queryHandlerTestEnv   *queryhandler.TestEnv
	commandHandlerTestEnv *commandhandler.TestEnv
}

func NewTestEnv() (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv: test.NewTestEnv("IntegrationTestEnv"),
	}

	os.Setenv("SUPERUSERS", "admin@monoskope.io")

	env.eventStoreTestEnv, err = eventstore.NewTestEnv()
	if err != nil {
		return nil, err
	}

	env.queryHandlerTestEnv, err = queryhandler.NewTestEnv(env.eventStoreTestEnv)
	if err != nil {
		return nil, err
	}

	env.commandHandlerTestEnv, err = commandhandler.NewTestEnv(env.eventStoreTestEnv, env.queryHandlerTestEnv)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) Shutdown() error {
	if err := env.queryHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.commandHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.eventStoreTestEnv.Shutdown(); err != nil {
		return err
	}

	return env.TestEnv.Shutdown()
}
