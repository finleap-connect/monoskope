package config

import (
	"time"

	"sigs.k8s.io/kind/pkg/errors"
)

var (
	ErrEmptyServer        = errors.New("has no server defined")
	ErrNoConfigExists     = errors.New("no valid monoconfig found")
	ErrAlreadyInitialized = errors.New("a configuartion already exists")
)

// Config holds the information needed to build connect to remote monoskope instance as a given user
type Config struct {
	// Server is the address of the Monoskope gateway (https://hostname:port).
	Server string `json:"server"`
	// AuthInformation contains information to authenticate against monoskope
	AuthInformation *AuthInformation `json:"auth-information,omitempty"`
}

type AuthInformation struct {
	// Token is the bearer token for authentication to the Monoskope gateway.
	Token        string    `json:"auth_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// NewConfig is a convenience function that returns a new Config object with defaults
func NewConfig() *Config {
	return &Config{}
}

// Validate validates if the config is valid
func (c *Config) Validate() error {
	if c.Server == "" {
		return ErrEmptyServer
	}
	return nil
}

func (c *Config) HasToken() bool {
	return c.AuthInformation != nil && c.AuthInformation.Token != ""
}
