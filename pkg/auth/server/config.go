package server

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"

type Config struct {
	auth.BaseConfig
	RootToken     *string
	ValidClientId string
}
