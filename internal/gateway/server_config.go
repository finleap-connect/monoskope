package gateway

import "gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"

type serverConfig struct {
	KeepAlive  bool
	AuthConfig *auth.Config
}

func NewServerConfig() *serverConfig {
	return &serverConfig{}
}
