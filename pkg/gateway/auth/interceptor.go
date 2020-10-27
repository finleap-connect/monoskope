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

// validates the authorization.
func (s *AuthInterceptor) valid(ctx context.Context, authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")

	// Perform the token validation here.
	_, err := s.authHandler.Verify(ctx, token)

	return err == nil
}

// ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func (s *AuthInterceptor) ensureValid(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpcutil.ErrMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !s.valid(ctx, md["authorization"]) {
		return grpcutil.ErrInvalidToken
	}
	return nil
}

func (s *AuthInterceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := s.ensureValid(ctx)
	if err != nil {
		return nil, err
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func (s *AuthInterceptor) StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := s.ensureValid(ss.Context())
	if err != nil {
		return err
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(srv, ss)
}
