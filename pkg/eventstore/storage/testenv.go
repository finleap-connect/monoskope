package storage

import (
	"github.com/go-pg/pg"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type EventStoreTestEnv struct {
	*test.TestEnv
	DB *pg.DB
}

func (env *EventStoreTestEnv) Shutdown() error {
	if env.DB != nil {
		defer env.DB.Close()
	}
	return env.TestEnv.Shutdown()
}
