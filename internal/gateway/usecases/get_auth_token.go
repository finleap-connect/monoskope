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

package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type getAuthTokenUsecase struct {
	*usecase.UseCaseBase
	request     *api.ClusterAuthTokenRequest
	result      *api.ClusterAuthTokenResponse
	signer      jwt.JWTSigner
	userRepo    repositories.ReadOnlyUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
	issuer      string
	validity    map[string]time.Duration
}

func NewGetAuthTokenUsecase(
	request *api.ClusterAuthTokenRequest,
	response *api.ClusterAuthTokenResponse,
	signer jwt.JWTSigner,
	userRepo repositories.ReadOnlyUserRepository,
	clusterRepo repositories.ReadOnlyClusterRepository,
	issuer string,
	validity map[string]time.Duration,
) usecase.UseCase {
	useCase := &getAuthTokenUsecase{
		usecase.NewUseCaseBase("get-auth-token"),
		request,
		response,
		signer,
		userRepo,
		clusterRepo,
		issuer,
		validity,
	}
	return useCase
}

func (s *getAuthTokenUsecase) Run(ctx context.Context) error {
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}
	userInfo := metadataManager.GetUserInformation()

	if userInfo.NotBefore.IsZero() {
		s.Log.Info("User's authentication time could not be determined", "id", userInfo.Id, "name", userInfo.Name, "email", userInfo.Email)
		return domainErrors.ErrUnauthorized
	}
	if time.Now().UTC().Sub(userInfo.NotBefore) > 1*time.Minute {
		s.Log.Info("User's authentication is too far in the past", "id", userInfo.Id, "name", userInfo.Name, "email", userInfo.Email)
		return domainErrors.ErrUnauthorized
	}

	s.Log.V(logger.DebugLevel).Info("Getting current user by id...", "id", userInfo.Id, "name", userInfo.Name, "email", userInfo.Email)
	user, err := s.userRepo.ByUserId(ctx, userInfo.Id)
	if err != nil {
		return err
	}

	clusterId := s.request.GetClusterId()
	s.Log.V(logger.DebugLevel).Info("Getting cluster by id...", "id", clusterId)
	cluster, err := s.clusterRepo.ByClusterId(ctx, clusterId)
	if err != nil {
		return err
	}

	k8sRole := s.request.GetRole()
	s.Log.V(logger.DebugLevel).Info("Validating role exists...", "role", k8sRole)
	if err := k8s.ValidateRole(k8sRole); err != nil {
		return err
	}

	username := strings.ToLower(strings.Split(user.GetEmail(), "@")[0])
	if s.request.GetRole() != string(k8s.DefaultRole) {
		username = fmt.Sprintf("%s-%s", username, s.request.GetRole())
	}

	s.Log.V(logger.DebugLevel).Info("Generating token for k8s user...", "username", username)
	token := auth.NewKubernetesAuthToken(&jwt.StandardClaims{
		Name:          user.GetName(),
		Email:         user.GetEmail(),
		EmailVerified: true,
	}, &jwt.ClusterClaim{
		ClusterId:       cluster.GetId(),
		ClusterName:     cluster.GetName(),
		ClusterUserName: username,
		ClusterRole:     s.request.Role,
	}, s.issuer, user.Id, s.validity[s.request.Role])
	s.Log.V(logger.DebugLevel).Info("Token issued successfully.", "RawToken", token, "Expiry", token.Expiry.Time().String())

	signedToken, err := s.signer.GenerateSignedToken(token)
	if err != nil {
		return err
	}
	s.Log.V(logger.DebugLevel).Info("Token signed successfully.", "SignedToken", signedToken)

	*s.result = api.ClusterAuthTokenResponse{
		AccessToken: signedToken,
		Expiry:      timestamppb.New(token.Expiry.Time()),
	}
	return nil
}
