package auth

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

// oauthAccessWithoutTLS supplies PerRPCCredentials from a given token.
// It does not need transport security configured since it is for use-cases where TLS is done by Linkerd for example.
type oauthAccessWithoutTLS struct {
	token oauth2.Token
}

// NewOauthAccessWithoutTLS constructs the PerRPCCredentials using a given token.
func NewOauthAccessWithoutTLS(token *oauth2.Token) credentials.PerRPCCredentials {
	return oauthAccessWithoutTLS{token: *token}
}

func (oa oauthAccessWithoutTLS) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": oa.token.Type() + " " + oa.token.AccessToken,
	}, nil
}

// RequireTransportSecurity returns false since the token source should be used where TLS is done by Linkerd for example.
func (oa oauthAccessWithoutTLS) RequireTransportSecurity() bool {
	return false
}
