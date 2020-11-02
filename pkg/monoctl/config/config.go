package config

import "sigs.k8s.io/kind/pkg/errors"

var (
	ErrEmptyServer    = errors.New("has no server defined")
	ErrNoConfigExists = errors.New("no valid monoconfig found")
)

// Config holds the information needed to build connect to remote monoskope instance as a given user
type Config struct {
	// Server is the address of the Monoskope gateway (https://hostname:port).
	Server string `json:"server"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate. This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `json:"insecure-skip-tls-verify,omitempty"`
	// CertificateAuthority is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityData contains PEM-encoded certificate authority certificates. Overrides CertificateAuthority
	// +optional
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
	// Token is the bearer token for authentication to the Monoskope gateway.
	// +optional
	Token string `json:"token,omitempty"`
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
