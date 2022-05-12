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

package grpc

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

// insecureOauthAccess supplies PerRPCCredentials from a given token.
type insecureOauthAccess struct {
	token oauth2.Token
}

// NewOauthAccessWithoutTransportSecurity constructs the PerRPCCredentials using a given token for test purposes where no TLS necessary.
func NewOauthAccessWithoutTransportSecurity(token *oauth2.Token) credentials.PerRPCCredentials {
	return insecureOauthAccess{token: *token}
}

func (oa insecureOauthAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": oa.token.Type() + " " + oa.token.AccessToken,
	}, nil
}

func (oa insecureOauthAccess) RequireTransportSecurity() bool {
	return false
}
