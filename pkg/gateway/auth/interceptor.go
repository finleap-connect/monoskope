package auth

import (
	"context"
	"strings"

	dexpb "github.com/dexidp/dex/api"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	// Authentication handler to verify bearer tokens
	authHandler *Handler
	// Logger interface
	log logger.Logger
}

func NewAuthInterceptor(dexClient dexpb.DexClient, authConfig *Config) (*AuthInterceptor, error) {
	authHandler, err := NewHandler(dexClient, authConfig)
	if err != nil {
		return nil, err
	}

	return &AuthInterceptor{
		log:         logger.WithName("auth"),
		authHandler: authHandler,
	}, nil
}

// valid validates the authorization.
func (s *AuthInterceptor) valid(ctx context.Context, authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	_, err := s.authHandler.Verify(ctx, token)
	return err == nil
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func (s *AuthInterceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpcutil.ErrMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !s.valid(ctx, md["authorization"]) {
		return nil, grpcutil.ErrInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
