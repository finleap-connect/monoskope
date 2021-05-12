package jwt

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	MonoskopeIssuer = "Monoskope"
	DefaultExpiry   = 24 * time.Hour
)

type ClusterBootstrapClaims struct {
	jwt.Claims
}

// Creates a new cluster bootstrap token
func NewClusterBootstrapClaims(subject string) *ClusterBootstrapClaims {
	now := time.Now().UTC()

	return &ClusterBootstrapClaims{
		Claims: jwt.Claims{
			ID:       uuid.New().String(),
			Issuer:   MonoskopeIssuer,
			Subject:  subject,
			Audience: jwt.Audience{subject},
			Expiry:   jwt.NewNumericDate(now.Add(DefaultExpiry)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
}
