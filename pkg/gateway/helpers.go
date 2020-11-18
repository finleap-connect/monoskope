package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateInsecureGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		opts = append(opts, grpc.WithPerRPCCredentials(auth.NewOauthAccessWithoutTLS(token)))
	}

	return createGatewayConnection(ctx, url, opts)
}

func CreateGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}

	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		opts = append(opts, grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)))
	}

	return createGatewayConnection(ctx, url, opts)
}

func createGatewayConnection(ctx context.Context, url string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithBlock())
	return grpc.DialContext(ctx, url, opts...)
}
