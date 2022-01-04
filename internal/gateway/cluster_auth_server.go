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
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/usecases"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
)

type clusterAuthApiServer struct {
	api.UnimplementedClusterAuthServer
	log         logger.Logger
	signer      jwt.JWTSigner
	userRepo    repositories.ReadOnlyUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
	issuer      string
	validity    map[string]time.Duration
}

func NewClusterAuthAPIServer(
	issuer string,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
	clusterRepo repositories.ReadOnlyClusterRepository,
	validity map[string]time.Duration,
) api.ClusterAuthServer {
	s := &clusterAuthApiServer{
		log:         logger.WithName("server"),
		signer:      signer,
		userRepo:    userRepo,
		clusterRepo: clusterRepo,
		issuer:      issuer,
		validity:    validity,
	}
	return s
}

func (s *clusterAuthApiServer) GetAuthToken(ctx context.Context, request *api.ClusterAuthTokenRequest) (*api.ClusterAuthTokenResponse, error) {
	response := new(api.ClusterAuthTokenResponse)
	uc := usecases.NewGetAuthTokenUsecase(request, response, s.signer, s.userRepo, s.clusterRepo, s.issuer, s.validity)
	err := uc.Run(ctx)
	if err != nil {
		return nil, err
	}
	return response, nil
}
