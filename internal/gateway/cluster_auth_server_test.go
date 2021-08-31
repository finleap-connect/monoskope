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

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/k8s"
)

var _ = Describe("Internal/Gateway/ClusterAuthServer", func() {
	ctx := context.Background()
	expectedRole := k8s.DefaultRole

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	It("can retrieve auth url", func() {
		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())

		clusters, err := env.ClusterRepo.GetAll(ctx, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(clusters)).To(BeNumerically(">=", 1))
		defer conn.Close()
		apiClient := api.NewClusterAuthClient(conn)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:     uuid.MustParse(env.AdminUser.GetId()),
			Name:   env.AdminUser.Name,
			Email:  env.AdminUser.Email,
			Issuer: "monoskope",
		})

		response, err := apiClient.GetAuthToken(mdManager.GetOutgoingGrpcContext(), &api.ClusterAuthTokenRequest{
			ClusterId: clusters[0].Id,
			Role:      string(expectedRole),
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
	})
})
