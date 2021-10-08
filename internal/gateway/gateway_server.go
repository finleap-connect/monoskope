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
	authConfig  *auth.Config
	authHandler *auth.Handler
	userRepo    repositories.ReadOnlyUserRepository
}

func NewGatewayAPIServer(authConfig *auth.Config, authHandler *auth.Handler, userRepo repositories.ReadOnlyUserRepository) api.GatewayServer {
	s := &gatewayApiServer{
		log:         logger.WithName("server"),
		authConfig:  authConfig,
		authHandler: authHandler,
		userRepo:    userRepo,
	}
	return s
}

func (s *gatewayApiServer) GetAuthInformation(ctx context.Context, state *api.AuthState) (*api.AuthInformation, error) {
	url, encodedState, err := s.authHandler.GetAuthCodeURL(state, s.authConfig.Scopes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid argument: %v", err)
	}

	return &api.AuthInformation{AuthCodeUrl: url, State: encodedState}, nil
}

func (s *gatewayApiServer) ExchangeAuthCode(ctx context.Context, code *api.AuthCode) (*api.AuthResponse, error) {
	// Exchange auth code with upstream identity provider
	s.log.Info("Exchanging auth code with issuer...")
	upstreamClaims, err := s.authHandler.Exchange(ctx, code.GetCode(), code.GetState(), code.CallbackUrl)
	if err != nil {
		s.log.Error(err, "User authentication failed.")
		return nil, err
	}
	s.log.Info("Exchanged successful, received upstream claims.", "name", upstreamClaims.Name, "email", upstreamClaims.Email)

	// Check that a user exists in monoskope
	s.log.Info("Checking user exists...", "email", upstreamClaims.Email)
	user, err := s.userRepo.ByEmail(ctx, upstreamClaims.Email)
	if err != nil {
		return nil, err
	}
	s.log.Info("User exists!", "name", user.Name, "email", user.Email, "id", user.ID())

	// Override upstream name
	upstreamClaims.Name = user.Name

	// Issue token
	signedToken, rawToken, err := s.authHandler.IssueToken(ctx, upstreamClaims, user.ID().String())
	if err != nil {
		s.log.Error(err, "Issuing token failed.")
		return nil, err
	}

	// Create response
	userInfo := &api.AuthResponse{
		AccessToken: signedToken,
		Expiry:      timestamppb.New(rawToken.Expiry.Time()),
		Username:    user.Name,
	}

	s.log.Info("User authenticated successfully.", "User", upstreamClaims.Email, "Expiry", userInfo.Expiry.AsTime().String())

	return userInfo, nil
}
