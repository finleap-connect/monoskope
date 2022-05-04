// Copyright 2022 Monoskope Authors
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
	"context"
	"time"

	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"gopkg.in/square/go-jose.v2"
)

type ServerConfig struct {
	URL           string
	TokenValidity time.Duration
}

// Server implements a very basic OIDC server which issues and validates tokens
type Server struct {
	config   *ServerConfig
	verifier jwt.JWTVerifier
	signer   jwt.JWTSigner
	log      logger.Logger
}

// NewServer creates a new OIDC server
func NewServer(config *ServerConfig, signer jwt.JWTSigner, verifier jwt.JWTVerifier) *Server {
	n := &Server{
		config:   config,
		signer:   signer,
		verifier: verifier,
		log:      logger.WithName("auth-server"),
	}
	return n
}

// IssueToken wraps the upstream claims in a JWT signed by Monoskope
func (n *Server) IssueToken(ctx context.Context, upstreamClaims *jwt.StandardClaims, userId string) (string, *jwt.AuthToken, error) {
	if upstreamClaims.FederatedClaims == nil {
		upstreamClaims.FederatedClaims = make(map[string]string)
	}

	token := NewAuthToken(upstreamClaims, n.config.URL, userId, n.config.TokenValidity)
	n.log.V(logger.DebugLevel).Info("Token issued successfully.", "RawToken", token, "Expiry", token.Expiry.Time().String())

	signedToken, err := n.signer.GenerateSignedToken(token)
	if err != nil {
		return "", nil, err
	}
	n.log.V(logger.DebugLevel).Info("Token signed successfully.", "SignedToken", signedToken)

	return signedToken, token, err
}

// Authorize parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (n *Server) Authorize(ctx context.Context, token string, claims interface{}) error {
	if err := n.verifier.Verify(token, claims); err != nil {
		return err
	}
	return nil
}

func (n *Server) Keys() *jose.JSONWebKeySet {
	return n.verifier.JWKS()
}
