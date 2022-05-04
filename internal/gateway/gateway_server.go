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

package gateway

import (
	"context"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gatewayApiServer struct {
	api.UnimplementedGatewayServer
	// Logger interface
	log logger.Logger
	//
	clientConfig *auth.ClientConfig
	authClient   *auth.Client
	authServer   *auth.Server
	userRepo     repositories.ReadOnlyUserRepository
}

func NewGatewayAPIServer(authConfig *auth.ClientConfig, client *auth.Client, server *auth.Server, userRepo repositories.ReadOnlyUserRepository) api.GatewayServer {
	s := &gatewayApiServer{
		log:          logger.WithName("server"),
		clientConfig: authConfig,
		authClient:   client,
		authServer:   server,
		userRepo:     userRepo,
	}
	return s
}

func (s *gatewayApiServer) RequestUpstreamAuthentication(ctx context.Context, request *api.UpstreamAuthenticationRequest) (*api.UpstreamAuthenticationResponse, error) {
	upstreamIDPUrl, encodedState, err := s.authClient.GetAuthCodeURL(request.GetCallbackUrl())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid argument: %v", err)
	}
	response := &api.UpstreamAuthenticationResponse{UpstreamIdpRedirect: upstreamIDPUrl, State: encodedState}
	s.log.V(logger.DebugLevel).Info("Upstream authentication requested.", "UpstreamAuthenticationResponse", response)
	return response, nil
}

func (s *gatewayApiServer) RequestAuthentication(ctx context.Context, request *api.AuthenticationRequest) (*api.AuthenticationResponse, error) {
	// Exchange auth code with upstream identity provider
	s.log.V(logger.DebugLevel).Info("Exchanging auth code with issuer...")
	upstreamClaims, err := s.authClient.Exchange(ctx, request.GetCode(), request.GetState())
	if err != nil {
		s.log.Error(err, "User authentication failed.")
		return nil, err
	}
	s.log.V(logger.DebugLevel).Info("Exchanged successful, received upstream claims.", "name", upstreamClaims.Name, "email", upstreamClaims.Email)

	// Check that a user exists in monoskope
	s.log.V(logger.DebugLevel).Info("Checking user exists...", "email", upstreamClaims.Email)
	user, err := s.userRepo.ByEmail(ctx, upstreamClaims.Email)
	if err != nil {
		return nil, err
	}
	s.log.V(logger.DebugLevel).Info("User exists!", "name", user.Name, "email", user.Email, "id", user.ID())

	// Override upstream name
	upstreamClaims.Name = user.Name

	// Issue token
	signedToken, rawToken, err := s.authServer.IssueToken(ctx, upstreamClaims, user.ID().String())
	if err != nil {
		s.log.Error(err, "Issuing token failed.", user.Name, "email", user.Email, "id", user.ID())
		return nil, err
	}

	// Create response
	response := &api.AuthenticationResponse{
		AccessToken: signedToken,
		Expiry:      timestamppb.New(rawToken.Expiry.Time()),
		Username:    user.Name,
	}

	s.log.Info("User authenticated successfully.", "User", upstreamClaims.Email, "Expiry", response.Expiry.AsTime().String())

	return response, nil
}
