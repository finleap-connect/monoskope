package eventstore

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// ServerConfig is the configuration for the API server
type serverConfig struct {
	KeepAlive bool
	Store     storage.Store
	Bus       messaging.EventBusPublisher
}

func NewServerConfig() *serverConfig {
	return &serverConfig{}
}
