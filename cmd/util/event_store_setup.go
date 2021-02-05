package util

import (
	"os"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
)

func NewEventStore() (eventsourcing.Store, error) {
	var dbUrl string

	if v := os.Getenv("DB_URL"); v != "" {
		dbUrl = v
	}

	conf, err := storage.NewPostgresStoreConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	err = conf.ConfigureTLS()
	if err != nil {
		return nil, err
	}

	return storage.NewPostgresEventStore(conf)
}
