package jwt

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	MonoskopeIssuer               = "Monoskope"
	ClusterBootstrapTokenValidity = 10 * time.Minute
)

type ClusterBootstrapToken struct {
	*jwt.Claims
}

// Creates a new cluster bootstrap token
func NewClusterBootstrapToken(subject string) *ClusterBootstrapToken {
	now := time.Now().UTC()

	return &ClusterBootstrapToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    MonoskopeIssuer,
			Subject:   subject,
			Audience:  jwt.Audience{subject},
			Expiry:    jwt.NewNumericDate(now.Add(ClusterBootstrapTokenValidity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

// IsValid returns if the token is not used too early or is expired
func (t *ClusterBootstrapToken) Validate() error {
	return t.ValidateWithLeeway(jwt.Expected{
		Issuer: MonoskopeIssuer,
	}, jwt.DefaultLeeway)
}
