package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateInsecureGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*grpc.ClientConn, error) {
	opts := make([]grpc.DialOption, 0)
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		opts = append(opts, grpc.WithPerRPCCredentials(auth.NewOauthAccessWithoutTLS(token)))
	}
	return grpcutil.CreateInsecureGrpcConnecton(ctx, url, opts)
}

func CreateGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		opts = append(opts, grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)))
	}
	return grpcutil.CreateGrpcConnection(ctx, url, opts)
}
