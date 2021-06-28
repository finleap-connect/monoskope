package usecases

import (
	"context"
	"fmt"
	"strings"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type getAuthTokenUsecase struct {
	*usecase.UseCaseBase
	request     *api.ClusterAuthTokenRequest
	result      *api.ClusterAuthTokenResponse
	signer      jwt.JWTSigner
	userRepo    repositories.ReadOnlyUserRepository
	clusterRepo repositories.ReadOnlyClusterRepository
}

func NewGetAuthTokenUsecase(request *api.ClusterAuthTokenRequest, result *api.ClusterAuthTokenResponse, signer jwt.JWTSigner, userRepo repositories.ReadOnlyUserRepository, clusterRepo repositories.ReadOnlyClusterRepository) usecase.UseCase {
	useCase := &getAuthTokenUsecase{
		UseCaseBase: usecase.NewUseCaseBase("get-auth-token"),
		request:     request,
		result:      result,
		signer:      signer,
		userRepo:    userRepo,
		clusterRepo: clusterRepo,
	}
	return useCase
}

func (s *getAuthTokenUsecase) Run(ctx context.Context) error {
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}
	userInfo := metadataManager.GetUserInformation()
	user, err := s.userRepo.ByUserId(ctx, userInfo.Id)
	if err != nil {
		return err
	}

	cluster, err := s.clusterRepo.ByClusterId(ctx, s.request.ClusterId)
	if err != nil {
		return err
	}

	if err := k8s.ValidateRole(s.request.GetRole()); err != nil {
		return err
	}

	username := strings.ToLower(user.Name)
	if s.request.GetRole() != string(k8s.DefaultRole) {
		username = fmt.Sprintf("%s-%s", username, s.request.GetRole())
	}

	token := jwt.NewKubernetesAuthToken(&jwt.StandardClaims{
		Name:          user.GetName(),
		Email:         user.GetEmail(),
		EmailVerified: true,
	}, &jwt.ClusterClaim{
		Id:       cluster.GetId(),
		Name:     cluster.GetName(),
		UserName: username,
		Role:     s.request.Role,
	}, user.Id, jwt.AuthTokenValidity)
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
