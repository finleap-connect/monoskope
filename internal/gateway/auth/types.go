package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

type Config struct {
	IdentityProvider     string
	IdentityProviderName string
	Scopes               []string
	OfflineAsScope       bool
	Nonce                string
	ClientId             string
	ClientSecret         string
	RedirectURIs         []string
}

func (conf *Config) String() string {
	return fmt.Sprintf(
		"IdentityProviderName: %s\\IdentityProvider: %s\\ņScopes: %v\\ņRedirectURIs: %v",
		conf.IdentityProviderName,
		conf.IdentityProvider,
		conf.Scopes,
		conf.RedirectURIs,
	)
}

type OpenIdConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksURL string `json:"jwks_uri"`
}

type State struct {
	Callback string `form:"callback" json:"callback,omitempty"`
}

func DecodeState(encoded string) (*State, error) {
	state := &State{}
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %v", err)
	}
	err = json.Unmarshal(data, state)
	return state, err
}

func (state *State) Encode() (string, error) {
	data, err := json.Marshal(state)
	if err != nil {
		return "", fmt.Errorf("error marshalling: %v", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	return encoded, nil
}

func (state *State) IsValid() bool {
	_, err := url.Parse(state.Callback)
	return err == nil
}
