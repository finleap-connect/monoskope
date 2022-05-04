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

package usecases

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/usecase"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type generateAPITokenUsecase struct {
	*usecase.UseCaseBase
	request  *api.APITokenRequest
	response *api.APITokenResponse
	signer   jwt.JWTSigner
	userRepo repositories.ReadOnlyUserRepository
	issuer   string
}

func NewGenerateAPITokenUsecase(
	request *api.APITokenRequest,
	response *api.APITokenResponse,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
	issuer string,
) usecase.UseCase {
	return &generateAPITokenUsecase{
		usecase.NewUseCaseBase("generate-api-token"),
		request,
		response,
		signer,
		userRepo,
		issuer,
	}
}

func (u *generateAPITokenUsecase) Run(ctx context.Context) error {
	// Validate scopes
	if len(u.request.GetAuthorizationScopes()) < 1 {
		return fmt.Errorf("At least one scope required.")
	}

	// Determine user
	var userId string
	standardClaims := new(jwt.StandardClaims)
	switch userRequest := u.request.User.(type) {
	case *api.APITokenRequest_UserId:
		userId = userRequest.UserId
		user, err := u.userRepo.ByUserId(ctx, uuid.MustParse(userRequest.UserId))
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
		standardClaims.Name = user.GetName()
		standardClaims.Email = user.GetEmail()
	case *api.APITokenRequest_Username:
		userId = userRequest.Username
		standardClaims.Name = userId
	default:
		return fmt.Errorf("user argument invalid")
	}

	// Generate and sign token
	token := auth.NewApiToken(standardClaims, u.issuer, userId, u.request.Validity.AsDuration(), u.request.GetAuthorizationScopes())
	u.Log.V(logger.DebugLevel).Info("Token issued successfully.", "RawToken", token, "Expiry", token.Expiry.Time().String())
	signedToken, err := u.signer.GenerateSignedToken(token)
	if err != nil {
		return err
	}
	u.Log.V(logger.DebugLevel).Info("Token signed successfully.", "SignedToken", signedToken)

	// Set response
	u.response.AccessToken = signedToken
	u.response.Expiry = timestamppb.New(token.Expiry.Time())

	return nil
}
