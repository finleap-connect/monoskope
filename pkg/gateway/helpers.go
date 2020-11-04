package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

func CreateGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		opts = append(opts, grpc.WithPerRPCCredentials(auth.NewOauthAccessWithoutTLS(token)))
	}
	opts = append(opts, grpc.WithBlock())
	return grpc.DialContext(ctx, url, opts...)
}
