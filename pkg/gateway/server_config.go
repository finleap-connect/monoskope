package gateway

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"

type ServerConfig struct {
	KeepAlive  bool
	AuthConfig *auth.Config
}
