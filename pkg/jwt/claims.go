package jwt

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	MonoskopeIssuer               = "Monoskope"
	AudienceMonoctl               = "monoctl"
	AudienceK8sAuth               = "k8sauth"
	AudienceM8Operator            = "m8operator"
	AuthTokenValidity             = 12 * time.Hour
	ClusterBootstrapTokenValidity = 10 * time.Minute
)

type StandardClaims struct {
	Name            string            `json:"name,omitempty"`           // User’s display name.
	Email           string            `json:"email,omitempty"`          // The email of the user.
	EmailVerified   bool              `json:"email_verified,omitempty"` // If the upstream provider has verified the email.
	Groups          []string          `json:"groups,omitempty"`         // A list of strings representing the groups a user is a member of.
	FederatedSub    string            `json:"sub"`                      // The connector ID and the user ID assigned to the user at the provider.
	FederatedClaims map[string]string `json:"federated_claims,omitempty"`
}

type ClusterClaim struct {
	Id       string `json:"id,omitempty"`        // Id of the cluster.
	Name     string `json:"name,omitempty"`      // Name of the cluster.
	UserName string `json:"user_name,omitempty"` // Name of the user in the cluster.
	Role     string `json:"role,omitempty"`      // Role the user has in the cluster.
}

type AuthToken struct {
	*jwt.Claims
	*StandardClaims
	ClusterClaim *ClusterClaim `json:"cluster_claims,omitempty"`
	ConnectorId  string        `json:"connector_id,omitempty"`
}

func NewAuthToken(claims *StandardClaims, userId, connectorId string) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    MonoskopeIssuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(AuthTokenValidity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Audience:  jwt.Audience{AudienceMonoctl},
		},
		StandardClaims: claims,
		ConnectorId:    connectorId,
	}
}

func NewKubernetesAuthToken(claims *StandardClaims, clusterClaim *ClusterClaim, userId string, validity time.Duration) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    MonoskopeIssuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(validity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Audience:  jwt.Audience{AudienceK8sAuth},
		},
		StandardClaims: claims,
		ClusterClaim:   clusterClaim,
		ConnectorId:    MonoskopeIssuer,
	}
}

func NewClusterBootstrapToken(claims *StandardClaims, userId, connectorId string) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    MonoskopeIssuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(ClusterBootstrapTokenValidity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Audience:  jwt.Audience{AudienceM8Operator},
		},
		StandardClaims: claims,
		ConnectorId:    connectorId,
	}
}

// IsValid returns if the token is not used too early or is expired
func (t *AuthToken) Validate(expectedAudience ...string) error {
	if len(expectedAudience) > 0 {
		for i := 0; i < len(expectedAudience); i++ {
			if t.validate(expectedAudience[i]) == nil {
				return nil
			}
		}
	}
	return t.validate(expectedAudience...)
}

func (t *AuthToken) validate(expectedAudience ...string) error {
	return t.ValidateWithLeeway(jwt.Expected{
		Issuer:   MonoskopeIssuer,
		Audience: expectedAudience,
		Time:     time.Now().UTC(),
	}, jwt.DefaultLeeway)
}
