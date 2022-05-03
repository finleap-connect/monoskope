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
