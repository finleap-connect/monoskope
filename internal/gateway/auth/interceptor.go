package auth

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AUTH_SCHEME = "bearer"
)

type AuthServerInterceptor struct {
	// Authentication handler to verify bearer tokens
	authHandler *Handler
}

func NewInterceptor(authHandler *Handler) (*AuthServerInterceptor, error) {
	return &AuthServerInterceptor{
		authHandler: authHandler,
	}, nil
}

// ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the
// handler.
func (s *AuthServerInterceptor) EnsureValid(ctx context.Context, fullMethodName string) (context.Context, error) {
	// Allow some methods to be called unauthenticated
	if strings.HasPrefix(fullMethodName, "/gateway.auth.Auth") || // allow all auth service methods without auth
		strings.HasPrefix(fullMethodName, "/grpc.health.v1.Health") { // allow health checks without auth
		return ctx, nil
	}

	token, err := grpc_auth.AuthFromMD(ctx, AUTH_SCHEME)
	if err != nil {
		return nil, err
	}

	// Perform the token validation here.
	claims, err := s.authHandler.Authorize(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	metadata.NewDomainMetadataManager(ctx).SetUserEmail(claims.Email)

	newCtx := context.WithValue(ctx, &ExtraClaims{}, claims)

	return newCtx, nil
}

// AuthFunc is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthFunc func(ctx context.Context, fullMethodName string) (context.Context, error)

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		newCtx, err = authFunc(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		newCtx, err = authFunc(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
