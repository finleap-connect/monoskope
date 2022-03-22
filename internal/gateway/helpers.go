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

package gateway

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/url"
	"strings"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"
)

func CreateInsecureConnection(ctx context.Context, url string) (*ggrpc.ClientConn, error) {
	return grpc.NewGrpcConnectionFactory(url).
		WithInsecure().
		WithRetry().
		WithBlock().
		Connect(ctx)
}

// defaultBearerTokenFromHeaders extracts the token from header
func defaultBearerTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers[auth.HeaderAuthorization]
	if !ok {
		return ""
	}
	split := strings.Split(authHeader, "Bearer ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

// clientCertificateFromHeaders extracts the client certificate from header
func clientCertificateFromHeaders(headers map[string]string) (*x509.Certificate, error) {
	pemData, ok := headers[auth.HeaderForwardedClientCert]
	if !ok || pemData == "" {
		return nil, errors.New("cert header is empty")
	}

	decodedValue, err := url.QueryUnescape(pemData)
	if err != nil {
		return nil, errors.New("could not unescape pem data from header")
	}

	block, _ := pem.Decode([]byte(decodedValue))
	if block == nil {
		return nil, errors.New("decoding pem failed")
	}

	return x509.ParseCertificate(block.Bytes)
}

// containsString returns if a given value is contained in a string slice
func containsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
