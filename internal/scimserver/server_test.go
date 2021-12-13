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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("internal/scimserver/Server", func() {
	var userId uuid.UUID

	getUsers := func() {
		req := httptest.NewRequest(http.MethodGet, "/Users", nil)
		rr := httptest.NewRecorder()
		testEnv.scimServer.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(rr.Body)
		Expect(err).To(Not(HaveOccurred()))
		testEnv.Log.Info(string(body))
	}

	createUser := func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/Users",
			strings.NewReader(`{"userName":"some.user@monoskope.io","schemas":["urn:scim:schemas:core:2.0"],"displayName":"Some User"}`),
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
	}

	getSpecificUser := func() {
		rr := httptest.NewRecorder()
		err := backoff.Retry(func() error {
			req := httptest.NewRequest(http.MethodGet, `/Users?filter=userName%20eq%20"some.user@monoskope.io"`, nil)
			testEnv.scimServer.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				return fmt.Errorf("wrong status code: %v", rr.Code)
			}
			return nil
		}, backoff.NewExponentialBackOff())
		Expect(err).To(Not(HaveOccurred()))

		body, err := ioutil.ReadAll(rr.Body)
		Expect(err).To(Not(HaveOccurred()))
		testEnv.Log.Info(string(body))
	}

	replaceUser := func() {
		var rr *httptest.ResponseRecorder
		err := backoff.Retry(func() error {
			rr = httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut, "/Users/"+userId.String(), strings.NewReader(`{"userName":"some.user@monoskope.io","schemas":["urn:scim:schemas:core:2.0"],"displayName":"Some User"}`),
			)
			testEnv.scimServer.ServeHTTP(rr, req)
			if rr.Code == http.StatusNotFound {
				return fmt.Errorf("wrong status code: %v", rr.Code)
			}
			return nil
		}, backoff.NewExponentialBackOff())
		Expect(err).To(Not(HaveOccurred()))
		Expect(rr.Code).To(Equal(http.StatusOK))
	}

	deleteUser := func() {
		req := httptest.NewRequest(http.MethodDelete, "/Users/"+userId.String(), nil)
		rr := httptest.NewRecorder()
		testEnv.scimServer.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusNoContent))
	}

	It("allows management of users", func() {
		By("querying users")
		getUsers()

		By("creating users")
		createUser()

		By("getting a specific user via filter")
		getSpecificUser()

		By("replacing a user")
		replaceUser()

		By("deleting a user")
		deleteUser()
	})

	getGroups := func() {
		req := httptest.NewRequest(http.MethodGet, "/Groups", nil)
		rr := httptest.NewRecorder()
		testEnv.scimServer.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(rr.Body)
		Expect(err).To(Not(HaveOccurred()))
		testEnv.Log.Info(string(body))
	}

	patchGroup := func() {
		req := httptest.NewRequest(http.MethodPatch, "/Groups/"+roles.IdFromRole(roles.Admin).String(), strings.NewReader(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":[{"value":[{"value":"%s"}],"op":"add","path":"members"}]}"`, userId.String())))
		rr := httptest.NewRecorder()
		testEnv.scimServer.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))

		body, err := ioutil.ReadAll(rr.Body)
		Expect(err).To(Not(HaveOccurred()))
		testEnv.Log.Info(string(body))
	}

	It("allows management of groups", func() {
		By("querying groups")
		getGroups()

		By("adding users to a group")
		createUser()
		patchGroup()
	})
})
