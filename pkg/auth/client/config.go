package client

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"

type Config struct {
	auth.BaseConfig
	Nonce        string
	ClientId     string
	ClientSecret string
	RedirectURI  string
}
