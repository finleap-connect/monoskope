package eventstore

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"

// ServerConfig is the configuration for the API server
type ServerConfig struct {
	KeepAlive bool
	Store     storage.Store
}
