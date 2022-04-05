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

package authn

import (
	"context"
	"fmt"

	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	mock_api "github.com/finleap-connect/monoskope/test/api/gateway"
	"github.com/golang/mock/gomock"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Test validation rules for cluster messages", func() {
	Context("Creating cluster", func() {
		var mockCtrl *gomock.Controller
		ctx := context.Background()

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		ctxWithToken := func(ctx context.Context, scheme string, token string) context.Context {
			md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
			nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
			return nCtx
		}

		It("should ensure rules are valid", func() {
			newCtx := ctxWithToken(ctx, "bearer", "onetoken")
			authzClient := mock_api.NewMockGatewayAuthZClient(mockCtrl)
			expectedMethodName := "test"

			middleware := NewAuthNMiddleware(authzClient).(*authNMiddleware)

			authzClient.EXPECT().Check(newCtx, &api.CheckRequest{
				FullMethodName: expectedMethodName,
			}).Return(&api.CheckResponse{
				Tags: []*api.CheckResponse_CheckResponseTag{
					{
						Key:   "test",
						Value: "test",
					},
				},
			}, nil)

			resultCtx, err := middleware.authNWithGateway(newCtx, expectedMethodName)
			Expect(err).ToNot(HaveOccurred())
			Expect(resultCtx).ToNot(BeNil())

			tags := grpc_ctxtags.Extract(resultCtx)
			Expect(tags).ToNot(BeNil())
			Expect(tags.Has("test")).To(BeTrue())
			Expect(tags.Values()["test"]).To(Equal("test"))
		})
	})

})
