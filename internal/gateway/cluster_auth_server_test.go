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
	"time"

	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/k8s"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Internal/Gateway/ClusterAuthServer", func() {
	ctx := context.Background()
	expectedRole := k8s.DefaultRole

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	It("can request a cluster auth token", func() {
		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())

		defer conn.Close()
		apiClient := api.NewClusterAuthClient(conn)

		mdManager.SetUserInformation(&metadata.UserInformation{
			Id:        uuid.MustParse(testEnv.TenantAdminUser.GetId()),
			Name:      testEnv.TenantAdminUser.Name,
			Email:     testEnv.TenantAdminUser.Email,
			NotBefore: time.Now().UTC(),
		})

		response, err := apiClient.GetAuthToken(mdManager.GetOutgoingGrpcContext(), &api.ClusterAuthTokenRequest{
			ClusterId: testEnv.TestClusterId.String(),
			Role:      string(expectedRole),
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
	})
})
