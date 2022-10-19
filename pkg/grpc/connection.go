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
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcConnectionFactory struct {
	opts []grpc.DialOption
	url  string
}

// NewGrpcConnectionFactory creates a new factory for gRPC connections.
func NewGrpcConnectionFactory(url string) *GrpcConnectionFactory {
	return &GrpcConnectionFactory{
		url:  url,
		opts: make([]grpc.DialOption, 0),
	}
}

// NewGrpcConnectionFactoryWithInsecure creates a new factory for gRPC connections and adds the following dial options: WithInsecure, WithBlock.
func NewGrpcConnectionFactoryWithInsecure(url string) *GrpcConnectionFactory {
	return NewGrpcConnectionFactory(url).
		WithInsecure().
		WithOpenTelemetry().
		WithBlock()
}

// WithInsecure adds a DialOption which disables transport security for this connection. Note that transport security is required unless WithInsecure is set.
func (factory *GrpcConnectionFactory) WithInsecure() *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return factory
}

// WithOSCaTransportCredentials adds a DialOption which configures a connection level security credentials (e.g., TLS/SSL) using the CAs known to the OS.
func (factory *GrpcConnectionFactory) WithOSCaTransportCredentials() *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	return factory
}

// WithPerRPCCredentials adds a DialOption which sets credentials and places auth state on each outbound RPC.
func (factory *GrpcConnectionFactory) WithPerRPCCredentials(creds credentials.PerRPCCredentials) *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithPerRPCCredentials(creds))
	return factory
}

// WithBlock adds a DialOption which makes caller of Dial blocks until the underlying connection is up. Without this, Dial returns immediately and connecting the server happens in background.
func (factory *GrpcConnectionFactory) WithBlock() *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithBlock())
	return factory
}

// WithOpenTelemetry adds a DialOption which adds OpenTelemetry to the client.
func (factory *GrpcConnectionFactory) WithOpenTelemetry() *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	factory.opts = append(factory.opts, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return factory
}

// WithRetry adds retrying with exponential backoff using the default retryable codes from grpc_retry.DefaultRetriableCodes.
func (factory *GrpcConnectionFactory) WithRetry() *GrpcConnectionFactory {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(10 * time.Millisecond)),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...),
		grpc_retry.WithMax(5),
	}

	factory.opts = append(factory.opts, grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)))
	factory.opts = append(factory.opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))

	return factory
}

// Connect creates a client connection based on the factory.
func (factory *GrpcConnectionFactory) WithTransportCredentials(creds credentials.TransportCredentials) *GrpcConnectionFactory {
	factory.opts = append(factory.opts, grpc.WithTransportCredentials(creds))
	return factory
}

// Connect creates a client connection based on the factory.
func (factory *GrpcConnectionFactory) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, factory.url, factory.opts...)
}

// ConnectWithTimeout creates a client connection based on the factory with a given timeout.
func (factory *GrpcConnectionFactory) ConnectWithTimeout(ctx context.Context, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return factory.Connect(ctx)
}

// NewClientWithAuthForward creates a new gRPC client which forwards the authentication bearer token received from the client (TLS optional)
func NewClientWithAuthForward[T any](ctx context.Context, addr string, requireTransportSecurity bool, clientFactory func(cc grpc.ClientConnInterface) T) (*grpc.ClientConn, T, error) {
	conn, err := NewGrpcConnectionFactory(addr).
		WithInsecure().
		WithPerRPCCredentials(NewForwardedOauthAccess(requireTransportSecurity)).
		WithBlock().
		ConnectWithTimeout(ctx, 10*time.Second)

	if err != nil {
		var result T
		return nil, result, err
	}

	return conn, clientFactory(conn), nil
}

// NewClientWithInsecureAuth (USE ONLY IF SECURED BY SERVICE MESH OR SIMILAR) creates a new gRPC client which sends the auth token without TLS
func NewClientWithInsecureAuth[T any](ctx context.Context, addr, authToken string, clientFactory func(cc grpc.ClientConnInterface) T) (*grpc.ClientConn, T, error) {
	conn, err := NewGrpcConnectionFactory(addr).
		WithInsecure().
		WithPerRPCCredentials(NewOauthAccessWithoutTransportSecurity(&oauth2.Token{AccessToken: authToken})).
		WithBlock().
		ConnectWithTimeout(ctx, 10*time.Second)

	if err != nil {
		var result T
		return nil, result, err
	}

	return conn, clientFactory(conn), nil
}

// NewClientWithInsecure creates a new gRPC client which connects without TLS
func NewClientWithInsecure[T any](ctx context.Context, addr string, clientFactory func(cc grpc.ClientConnInterface) T) (*grpc.ClientConn, T, error) {
	conn, err := NewGrpcConnectionFactoryWithInsecure(addr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		var result T
		return nil, result, errors.TranslateToGrpcError(err)
	}

	return conn, clientFactory(conn), nil
}
