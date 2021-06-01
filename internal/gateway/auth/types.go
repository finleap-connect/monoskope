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

// TokenReviewRequest is issued by K8s to this service
type TokenReviewRequest struct {
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Spec       TokenReviewSpec   `json:"spec"`
	Status     TokenReviewStatus `json:"status"`
}

// TokenReviewSpec contains the token being reviewed
type TokenReviewSpec struct {
	Token string `json:"token"`
}

// TokenReviewStatus is the result of the token authentication request.
type TokenReviewStatus struct {
	// Authenticated is true if the token is valid
	Authenticated bool `json:"authenticated,omitempty"`
	// User contains information about the authenticated user.
	User UserInfo `json:"user,omitempty"`
}

// UserInfo contains information about the user
type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	Username string `json:"username,omitempty"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	UID string `json:"uid,omitempty"`
	// The names of groups this user is a part of.
	Groups []string `json:"groups,omitempty"`
	// Any additional information provided by the authenticator.
	Extra map[string][]string `json:"extra,omitempty"`
}
