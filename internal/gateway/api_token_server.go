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
	"fmt"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type apiTokenServer struct {
	api.UnimplementedAPITokenServer
	log      logger.Logger
	signer   jwt.JWTSigner
	userRepo repositories.ReadOnlyUserRepository
	url      string
}

func NewAPITokenServer(
	url string,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
) api.APITokenServer {
	s := &apiTokenServer{
		log:      logger.WithName("server"),
		signer:   signer,
		userRepo: userRepo,
		url:      url,
	}
	return s
}

func (s *apiTokenServer) RequestAPIToken(ctx context.Context, request *api.APITokenRequest) (*api.APITokenResponse, error) {
	if len(request.GetAuthorizationScopes()) < 1 {
		return nil, fmt.Errorf("At least one scope required.")
	}

	standardClaims := new(jwt.StandardClaims)

	var userId string
	switch u := request.User.(type) {
	case *api.APITokenRequest_UserId:
		userId = u.UserId
		user, err := s.userRepo.ByUserId(ctx, uuid.MustParse(u.UserId))
		if err != nil {
			return nil, errors.TranslateToGrpcError(err)
		}
		standardClaims.Name = user.GetName()
		standardClaims.Email = user.GetEmail()
	case *api.APITokenRequest_Username:
		userId = u.Username
	default:
		return nil, fmt.Errorf("user argument invalid")
	}

	token := auth.NewApiToken(standardClaims, s.url, userId, request.Validity.AsDuration(), request.GetAuthorizationScopes())

	signedToken, err := s.signer.GenerateSignedToken(token)
	if err != nil {
		return nil, err
	}

	result := new(api.APITokenResponse)
	result.AccessToken = signedToken
	result.Expiry = timestamppb.New(token.Expiry.Time())
	return result, nil
}
