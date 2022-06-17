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
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	jose_jwt "gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("internal/gateway/auth/token", func() {
	expectedIssuer := "https://localhost"
	expectedValidity := time.Hour * 1
	It("validate cluster bootstrap token", func() {
		t := NewClusterBootstrapToken(&jwt.StandardClaims{}, expectedIssuer, "me")
		Expect(t.Validate(expectedIssuer)).ToNot(HaveOccurred())
	})
	It("fail validate auth token", func() {
		t := NewAuthToken(&jwt.StandardClaims{}, expectedIssuer, "me", expectedValidity)
		t.Expiry = jose_jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * -12))
		Expect(t.Validate(expectedIssuer)).To(HaveOccurred())
	})

	It("validate api token", func() {
		t := NewApiToken(&jwt.StandardClaims{}, expectedIssuer, "me", expectedValidity, []gateway.AuthorizationScope{gateway.AuthorizationScope_WRITE_SCIM})
		t.Expiry = jose_jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 1))
		Expect(t.Validate(expectedIssuer)).ToNot(HaveOccurred())
		Expect(t.Scope).To(Equal(gateway.AuthorizationScope_WRITE_SCIM.String()))
	})
})
