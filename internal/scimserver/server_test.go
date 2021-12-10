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

package scimserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("internal/scimserver/Server", func() {
	var userId uuid.UUID

	When("getting users", func() {
		It("returns the list of users and StatusOK", func() {
			req := httptest.NewRequest(http.MethodGet, "/Users", nil)
			rr := httptest.NewRecorder()
			testEnv.scimServer.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(rr.Body)
			Expect(err).To(Not(HaveOccurred()))
			testEnv.Log.Info(string(body))
		})
	})
	When("creating user", func() {
		It("returns the newly created user and StatusCreated", func() {
			req := httptest.NewRequest(
				http.MethodPost,
				"/Users",
				strings.NewReader(`{"userName":"some.user@monoskope.io","schemas":["urn:scim:schemas:core:1.0"],"displayName":"Some User"}`),
			)
			rr := httptest.NewRecorder()
			testEnv.scimServer.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusCreated))

			body, err := ioutil.ReadAll(rr.Body)
			Expect(err).To(Not(HaveOccurred()))
			Expect(body).To(MatchRegexp(`^{"displayName":"Some User","id":"[0-9a-z\-]+","meta":{"resourceType":"User","location":"Users/[0-9a-z\-]+"},"schemas":\["urn:ietf:params:scim:schemas:core:2.0:User"\],"userName":"some.user@monoskope.io"}$`))
			testEnv.Log.Info(string(body))

			bodyMap := make(map[string]interface{})
			err = json.Unmarshal(body, &bodyMap)
			Expect(err).To(Not(HaveOccurred()))
			userIdStr := bodyMap["id"].(string)
			userId, err = uuid.Parse(userIdStr)
			Expect(err).To(Not(HaveOccurred()))
		})
	})
	When("deleting user", func() {
		It("returns status StatusNoContent", func() {
			req := httptest.NewRequest(http.MethodDelete, "/Users/"+userId.String(), nil)
			rr := httptest.NewRecorder()
			testEnv.scimServer.ServeHTTP(rr, req)
			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})
})
