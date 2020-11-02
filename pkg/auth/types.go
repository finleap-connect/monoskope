package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type BaseConfig struct {
	IssuerURL      string
	OfflineAsScope bool
}

type ExtraClaims struct {
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Groups        []string `json:"groups"`
}

type State struct {
	Callback string `form:"callback" json:"callback,omitempty"`
}

type AuthCodeURLConfig struct {
	Scopes        []string
	Clients       []string
	OfflineAccess bool
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
