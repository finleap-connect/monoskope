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

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/usecases"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type clusterAuthApiServer struct {
	api.UnimplementedClusterAuthServer
	// Logger interface
	log logger.Logger
	//
	signer      jwt.JWTSigner
	userRepo    repositories.ReadOnlyUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
	url         string
	validity    time.Duration
}

func NewClusterAuthAPIServer(
	url string,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
	clusterRepo repositories.ReadOnlyClusterRepository,
	validity time.Duration,
) api.ClusterAuthServer {
	s := &clusterAuthApiServer{
		log:         logger.WithName("server"),
		signer:      signer,
		userRepo:    userRepo,
		clusterRepo: clusterRepo,
		url:         url,
		validity:    validity,
	}
	return s
}

func (s *clusterAuthApiServer) GetAuthToken(ctx context.Context, request *api.ClusterAuthTokenRequest) (*api.ClusterAuthTokenResponse, error) {
	result := new(api.ClusterAuthTokenResponse)
	uc := usecases.NewGetAuthTokenUsecase(request, result, s.signer, s.userRepo, s.clusterRepo, s.url, s.validity)
	err := uc.Run(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
