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

package auth

import (
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	jose_jwt "gopkg.in/square/go-jose.v2/jwt"
)

func NewAuthToken(claims *jwt.StandardClaims, issuer, userId string, validity time.Duration) *jwt.AuthToken {
	now := time.Now().UTC()

	return &jwt.AuthToken{
		Claims: &jose_jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jose_jwt.NewNumericDate(now.Add(validity)),
			NotBefore: jose_jwt.NewNumericDate(now),
			IssuedAt:  jose_jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
		Scope:          gateway.AuthorizationScope_API.String(),
	}
}

func NewKubernetesAuthToken(claims *jwt.StandardClaims, clusterClaim *jwt.ClusterClaim, issuer, userId string, validity time.Duration) *jwt.AuthToken {
	now := time.Now().UTC()

	return &jwt.AuthToken{
		Claims: &jose_jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jose_jwt.NewNumericDate(now.Add(validity)),
			NotBefore: jose_jwt.NewNumericDate(now),
			IssuedAt:  jose_jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
		ClusterClaim:   clusterClaim,
		Scope:          gateway.AuthorizationScope_NONE.String(),
	}
}

const (
	ClusterBootstrapTokenValidity = 10 * time.Minute
)

func NewClusterBootstrapToken(claims *jwt.StandardClaims, issuer, userId string) *jwt.AuthToken {
	now := time.Now().UTC()

	return &jwt.AuthToken{
		Claims: &jose_jwt.Claims{
			ID:        uuid.New().String(),
			Issuer:    issuer,
			Subject:   userId,
			Expiry:    jose_jwt.NewNumericDate(now.Add(ClusterBootstrapTokenValidity)),
			NotBefore: jose_jwt.NewNumericDate(now),
			IssuedAt:  jose_jwt.NewNumericDate(now),
		},
		StandardClaims: claims,
		Scope:          gateway.AuthorizationScope_WRITE_K8SOPERATOR.String(),
	}
}
