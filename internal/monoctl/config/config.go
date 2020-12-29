package config

import (
	"errors"
	"time"
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

// HasToken checks if the the config contains AuthInformation
func (c *Config) HasAuthInformation() bool {
	return c.AuthInformation != nil
}

type AuthInformation struct {
	// Token is the bearer token for authentication to the Monoskope gateway.
	Token        string    `json:"auth_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// IsValid checks that Token is not empty and is not expired
func (a *AuthInformation) IsValid() bool {
	return a.HasToken() && !a.IsTokenExpired()
}

// HasToken checks that Token is not empty
func (a *AuthInformation) HasToken() bool {
	return a.Token != ""
}

// HasRefreshToken checks that RefreshToken is not empty
func (a *AuthInformation) HasRefreshToken() bool {
	return a.RefreshToken != ""
}

// IsTokenExpired checks if the auth token is expired
func (a *AuthInformation) IsTokenExpired() bool {
	return a.Expiry.Before(time.Now().UTC().Add(5 * time.Minute)) // check if token is valid for at least five more minutes
}
