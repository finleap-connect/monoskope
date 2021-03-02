package config

import (
	"errors"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	ErrEmptyServer        = errors.New("has no server defined")
	ErrNoConfigExists     = errors.New("no valid monoconfig found")
	ErrAlreadyInitialized = errors.New("a configuartion already exists")
)

// Config holds the information needed to build connect to remote monoskope instance as a given user
type Config struct {
	// Server is the address of the Monoskope gateway (https://hostname:port).
	Server string `yaml:"server"`
	// AuthInformation contains information to authenticate against monoskope
	AuthInformation *AuthInformation `yaml:"auth-information,omitempty"`
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

func (c *Config) String() (string, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// HasToken checks if the the config contains AuthInformation
func (c *Config) HasAuthInformation() bool {
	return c.AuthInformation != nil
}

type AuthInformation struct {
	AccessToken  string     `yaml:"auth_token,omitempty"`
	IdToken      string     `yaml:"id_token,omitempty"`
	RefreshToken string     `yaml:"refresh_token,omitempty"`
	Expiry       *time.Time `yaml:"expiry,omitempty"`
}

// IsValid checks that Token is not empty and is not expired
func (a *AuthInformation) IsValid() bool {
	return a.HasToken() && !a.IsTokenExpired()
}

// HasToken checks that Token is not empty
func (a *AuthInformation) HasToken() bool {
	return a.AccessToken != "" || a.IdToken != ""
}

func (a *AuthInformation) GetToken() string {
	if a.IdToken != "" {
		return a.IdToken
	}
	if a.AccessToken != "" {
		return a.AccessToken
	}
	return ""
}

// HasRefreshToken checks that RefreshToken is not empty
func (a *AuthInformation) HasRefreshToken() bool {
	return a.RefreshToken != ""
}

// IsTokenExpired checks if the auth token is expired
func (a *AuthInformation) IsTokenExpired() bool {
	return a.Expiry != nil && !a.Expiry.IsZero() && a.Expiry.Before(time.Now().UTC().Add(5*time.Minute)) // check if token is valid for at least five more minutes
}
