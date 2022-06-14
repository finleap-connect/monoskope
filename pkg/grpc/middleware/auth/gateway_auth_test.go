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

package auth

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	mock_api "github.com/finleap-connect/monoskope/internal/test/api/gateway"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	domain_metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Gateway Auth Middleware", func() {
	Context("Auth requests", func() {
		var mockCtrl *gomock.Controller
		ctx := context.Background()

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		ctxWithToken := func(ctx context.Context, scheme string, token string) context.Context {
			md := metadata.Pairs(auth.HeaderAuthorization, fmt.Sprintf("%s %v", scheme, token))
			nCtx := metautils.NiceMD(md).ToIncoming(ctx)
			return nCtx
		}

		It("should ensure user is authenticated and authorized", func() {
			expectedToken := "onetoken"
			newCtx := ctxWithToken(ctx, "bearer", expectedToken)
			authClient := mock_api.NewMockGatewayAuthClient(mockCtrl)
			expectedMethodName := "test"
			expectedUserId := uuid.New()
			expectedRole := roles.Admin
			exptectedScope := scopes.Tenant
			expectedResource := uuid.New()

			command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding)
			_, err := cmd.AddCommandData(command,
				&cmdData.CreateUserRoleBindingCommandData{
					UserId:   expectedUserId.String(),
					Role:     string(expectedRole),
					Scope:    string(exptectedScope),
					Resource: wrapperspb.String(expectedResource.String()),
				},
			)
			Expect(err).ToNot(HaveOccurred())

			bytes, err := protojson.Marshal(command)
			Expect(err).ToNot(HaveOccurred())

			middleware := NewAuthMiddleware(authClient, nil).(*authMiddleware)

			authClient.EXPECT().Check(newCtx, &api.CheckRequest{
				FullMethodName: expectedMethodName,
				AccessToken:    expectedToken,
				Request:        bytes,
			}, gomock.Any()).Return(&api.CheckResponse{
				Tags: []*api.CheckResponse_CheckResponseTag{
					{
						Key:   auth.HeaderAuthId,
						Value: expectedUserId.String(),
					},
				},
			}, nil)

			resultCtx, err := middleware.authWithGateway(newCtx, expectedMethodName, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(resultCtx).ToNot(BeNil())

			m, err := domain_metadata.NewDomainMetadataManager(resultCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(m).ToNot(BeNil())

			userInfo := m.GetUserInformation()
			Expect(userInfo.Id).To(Equal(expectedUserId))
		})
	})

})
