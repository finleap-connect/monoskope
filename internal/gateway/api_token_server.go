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

	"github.com/finleap-connect/monoskope/internal/gateway/usecases"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
)

type apiTokenServer struct {
	api.UnimplementedAPITokenServer
	log      logger.Logger
	signer   jwt.JWTSigner
	userRepo repositories.ReadOnlyUserRepository
	issuer   string
}

func NewAPITokenServer(
	issuer string,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
) api.APITokenServer {
	s := &apiTokenServer{
		log:      logger.WithName("server"),
		signer:   signer,
		userRepo: userRepo,
		issuer:   issuer,
	}
	return s
}

func (s *apiTokenServer) RequestAPIToken(ctx context.Context, request *api.APITokenRequest) (*api.APITokenResponse, error) {
	response := new(api.APITokenResponse)
	uc := usecases.NewGenerateAPITokenUsecase(request, response, s.signer, s.userRepo, s.issuer)
	err := uc.Run(ctx)
	if err != nil {
		return nil, err
	}
	return response, nil
}
