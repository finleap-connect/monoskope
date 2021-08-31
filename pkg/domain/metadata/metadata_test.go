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

package metadata

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
)

var _ = Describe("Managing Metadata", func() {
	It("should set user information", func() {

		expectedUserId := uuid.New()

		ctx := context.Background()
		mdManager, err := NewDomainMetadataManager(ctx)
		Expect(err).ToNot(HaveOccurred())

		mdManager.SetUserInformation(&UserInformation{
			Id:     expectedUserId,
			Name:   "admin",
			Email:  "admin@monoskope.io",
			Issuer: "monoskope",
		})

		Expect(mdManager.GetMetadata()[auth.HeaderAuthId]).To(Equal(expectedUserId.String()))

	})
})
