package eventstore

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"

type ServerConfig struct {
	KeepAlive bool
	Store     storage.Store
}
