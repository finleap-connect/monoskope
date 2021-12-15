// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	ClusterBootstrapTokenValidity = 10 * time.Minute
)

// https://www.iana.org/assignments/jwt/jwt.xhtml
type StandardClaims struct {
	Name            string            `json:"name,omitempty"`             // User’s display name.
	Email           string            `json:"email,omitempty"`            // The email of the user.
	EmailVerified   bool              `json:"email_verified,omitempty"`   // If the upstream provider has verified the email.
	Groups          []string          `json:"groups,omitempty"`           // A list of strings representing the groups a user is a member of.
	FederatedClaims map[string]string `json:"federated_claims,omitempty"` // Claims from any upstream IDP.
	Scope           string            `json:"scope"`                      // Space-separated list of scopes associated with the token.
}

type ClusterClaim struct {
	ClusterId       string `json:"cluster_id,omitempty"`       // Id of the cluster.
	ClusterName     string `json:"cluster_name,omitempty"`     // Name of the cluster.
	ClusterUserName string `json:"cluster_username,omitempty"` // Name of the user in the cluster.
	ClusterRole     string `json:"cluster_role,omitempty"`     // Role the user has in the cluster.
}

type AuthToken struct {
	*jwt.Claims
	*StandardClaims
	*ClusterClaim
}

func NewAuthToken(claims *StandardClaims, issuer, userId string, validity time.Duration) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(validity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
	}
}

func NewKubernetesAuthToken(claims *StandardClaims, clusterClaim *ClusterClaim, issuer, userId string, validity time.Duration) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(validity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
		ClusterClaim:   clusterClaim,
	}
}

func NewClusterBootstrapToken(claims *StandardClaims, issuer, userId string) *AuthToken {
	now := time.Now().UTC()

	return &AuthToken{
		Claims: &jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jwt.NewNumericDate(now.Add(ClusterBootstrapTokenValidity)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
	}
}

// IsValid returns if the token is not used too early or is expired
func (t *AuthToken) Validate(issuer string, expectedAudience ...string) error {
	if len(expectedAudience) > 0 {
		for i := 0; i < len(expectedAudience); i++ {
			if t.validate(issuer, expectedAudience[i]) == nil {
				return nil
			}
		}
	}
	return t.validate(issuer, expectedAudience...)
}

func (t *AuthToken) validate(issuer string, expectedAudience ...string) error {
	return t.ValidateWithLeeway(jwt.Expected{
		Issuer:   issuer,
		Audience: expectedAudience,
		Time:     time.Now().UTC(),
	}, jwt.DefaultLeeway)
}
