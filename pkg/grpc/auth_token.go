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

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/credentials"
)

// oauthAccess supplies PerRPCCredentials from a given token.
type oauthAccess struct {
	requireTransportSecurity bool
}

// NewForwardedOauthAccess constructs the PerRPCCredentials which forwards the authorization from header
func NewForwardedOauthAccess(requireTransportSecurity bool) credentials.PerRPCCredentials {
	return oauthAccess{requireTransportSecurity}
}

func (oa oauthAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	// Forward authorization header
	token, err := grpc_auth.AuthFromMD(ctx, auth.AuthScheme)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"authorization": auth.AuthScheme + " " + token,
	}, nil
}

func (oa oauthAccess) RequireTransportSecurity() bool {
	return oa.requireTransportSecurity
}
