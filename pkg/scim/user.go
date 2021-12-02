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

package scim

import (
	"encoding/json"

	"github.com/elimity-com/scim"
)

type emailAttribute struct {
	Primary bool   `json:"primary"`
	Value   string `json:"value"`
}

type userAttributes struct {
	UserName string           `json:"userName"`
	Active   bool             `json:"active"`
	Emails   []emailAttribute `json:"emails"`
}

// GetPrimaryMail returns the primary email address with fallback to the first one found
func (u *userAttributes) GetPrimaryMail() string {
	var firstMail string
	for _, email := range u.Emails {
		if email.Primary {
			return email.Value
		} else if firstMail == "" {
			firstMail = email.Value
		}
	}
	return firstMail
}

// NewUserAttribute converts the SCIM resource attributes given to an instance of the userResource struct
func NewUserAttribute(attributes scim.ResourceAttributes) (*userAttributes, error) {
	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	userResource := new(userAttributes)
	err = json.Unmarshal(attributesJson, userResource)
	if err != nil {
		return nil, err
	}
	return userResource, nil
}
