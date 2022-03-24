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

package authn

import (
	"context"

	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	mgrpc "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/grpc/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authNMiddleware struct {
	factory *mgrpc.GrpcConnectionFactory
}

func NewAuthNMiddleware(authnServiceURL string) middleware.GRPCMiddleware {
	return &authNMiddleware{
		mgrpc.NewGrpcConnectionFactory(authnServiceURL).WithOSCaTransportCredentials().WithRetry().WithBlock(),
	}
}

// authnWithGateway calls the Gateway to authenticate the request and enriches the new context with tags set by the Gateway.
func (m *authNMiddleware) authnWithGateway(ctx context.Context, fullMethodName string) (context.Context, error) {
	// Connect to Gateway
	conn, err := m.factory.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewGatewayAuthZClient(conn)

	// Check request is authenticated and authorized
	response, err := client.Check(ctx, &api.CheckRequest{
		FullMethodName: fullMethodName,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "gateway authN failed: %v", err)
	}

	// Add tags from response to context
	tags := grpc_ctxtags.Extract(ctx)
	for _, tag := range response.GetTags() {
		tags.Set(tag.Key, tag.Value)
	}
	newCtx := grpc_ctxtags.SetInContext(ctx, tags)

	// Return new context with auth infos
	return newCtx, nil
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func (m *authNMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error

		newCtx, err = m.authnWithGateway(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func (m *authNMiddleware) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		newCtx, err = m.authnWithGateway(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
