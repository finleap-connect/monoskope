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
	"github.com/elimity-com/scim/optional"
	. "github.com/elimity-com/scim/schema"
)

const (
	UserNameAttribute    = "userName"
	DisplayNameAttribute = "displayName"
	ActiveAttribute      = "active"
)

// MonoskopeUserSchema returns the default "User" Resource Schema.
func MonoskopeUserSchema() Schema {
	return Schema{
		Description: optional.NewString("Monoskope User Account"),
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        optional.NewString("User"),
		Attributes: []CoreAttribute{
			SimpleCoreAttribute(SimpleStringParams(StringParams{
				Description: optional.NewString("Unique identifier for the User, used by the user to directly authenticate to the service provider (email address). Each User MUST include a non-empty userName value. This identifier MUST be unique across the service provider's entire set of Users. REQUIRED."),
				Name:        "userName",
				Required:    true,
				Uniqueness:  AttributeUniquenessServer(),
			})),
			SimpleCoreAttribute(SimpleStringParams(StringParams{
				Description: optional.NewString("The name of the User, suitable for display to end-users. The name SHOULD be the full name of the User being described, if known."),
				Name:        "displayName",
			})),
			SimpleCoreAttribute(SimpleBooleanParams(BooleanParams{
				Description: optional.NewString("A Boolean value indicating the User's administrative status."),
				Name:        "active",
			})),
		},
	}
}
