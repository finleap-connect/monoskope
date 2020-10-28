package server

import (
	"context"

	dexpb "github.com/dexidp/dex/api"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerInterceptor struct {
	// Authentication handler to verify bearer tokens
	authHandler *Handler
}

func NewInterceptor(dexClient dexpb.DexClient, authConfig *auth.Config) (*AuthServerInterceptor, error) {
	authHandler, err := NewHandler(dexClient, authConfig)
	if err != nil {
		return nil, err
	}

	return &AuthServerInterceptor{
		authHandler: authHandler,
	}, nil
}

// ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func (s *AuthServerInterceptor) EnsureValid(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	// Perform the token validation here.
	claims, err := s.authHandler.Authorize(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	grpc_ctxtags.Extract(ctx).Set("auth.sub", claims.Email)

	newCtx := context.WithValue(ctx, &auth.ExtraClaims{}, claims)
	return newCtx, nil
}
