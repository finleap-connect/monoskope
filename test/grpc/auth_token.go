package grpc

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

// oauthAccess supplies PerRPCCredentials from a given token.
type oauthAccess struct {
	token oauth2.Token
}

// NewOauthAccessWithoutTransportSecurity constructs the PerRPCCredentials using a given token for test purposes where no TLS necessary.
func NewOauthAccessWithoutTransportSecurity(token *oauth2.Token) credentials.PerRPCCredentials {
	return oauthAccess{token: *token}
}

func (oa oauthAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": oa.token.Type() + " " + oa.token.AccessToken,
	}, nil
}

func (oa oauthAccess) RequireTransportSecurity() bool {
	return false
}
