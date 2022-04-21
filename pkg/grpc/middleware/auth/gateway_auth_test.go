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

package auth

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	mock_api "github.com/finleap-connect/monoskope/test/api/gateway"
	"github.com/golang/mock/gomock"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
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

			middleware := NewAuthMiddleware(authClient, nil).(*authMiddleware)

			authClient.EXPECT().Check(newCtx, &api.CheckRequest{
				FullMethodName: expectedMethodName,
				AccessToken:    expectedToken,
			}, gomock.Any()).Return(&api.CheckResponse{
				Tags: []*api.CheckResponse_CheckResponseTag{
					{
						Key:   "test",
						Value: "test",
					},
				},
			}, nil)

			resultCtx, err := middleware.authWithGateway(newCtx, expectedMethodName, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(resultCtx).ToNot(BeNil())

			tags := grpc_ctxtags.Extract(resultCtx)
			Expect(tags).ToNot(BeNil())
			Expect(tags.Has("test")).To(BeTrue())
			Expect(tags.Values()["test"]).To(Equal("test"))
		})
	})

})
