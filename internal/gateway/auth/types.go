package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

type Config struct {
	IssuerIdentifier string
	IssuerURL        string
	OfflineAsScope   bool
	Nonce            string
	ClientId         string
	ClientSecret     string
	Scopes           []string
	RedirectURIs     []string
}

func (conf *Config) String() string {
	return fmt.Sprintf(
		"IssuerIdentifier: %s\\nIssuerURL: %s\\ņScopes: %v\\ņRedirectURIs: %v",
		conf.IssuerIdentifier,
		conf.IssuerURL,
		conf.Scopes,
		conf.RedirectURIs,
	)
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
