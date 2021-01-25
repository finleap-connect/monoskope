package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
)

type Config struct {
	IssuerURL      string
	OfflineAsScope bool
	Nonce          string
	ClientId       string
	ClientSecret   string
}

type Claims struct {
	Subject       string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Issuer        string `json:"iss"`
}

// Converts the claims provided by the IDP to proto
func (c *Claims) ToProto() *api_common.UserMetadata {
	return &api_common.UserMetadata{
		Subject: c.Subject,
		Issuer:  c.Issuer,
		Email:   c.Email,
	}
}

type AuthCodeURLConfig struct {
	Scopes        []string
	OfflineAccess bool
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
