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

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Test validation rules for cluster messages", func() {
	Context("Creating cluster", func() {
		ctxWithToken := func(ctx context.Context, scheme string, token string) context.Context {
			md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
			nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
			return nCtx
		}

		It("should ensure rules are valid", func() {
			ctx := context.Background()
			newCtx := ctxWithToken(ctx, "bearer", "onetoken")
			middleware := NewAuthNMiddleware("dummyurl").(*authNMiddleware)
			resultCtx, err := middleware.authnWithGateway(newCtx, "somemethod")
			Expect(err).ToNot(HaveOccurred())
			Expect(resultCtx).ToNot(BeNil())
		})
	})

})
