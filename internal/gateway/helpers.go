package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateInsecureGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*ggrpc.ClientConn, error) {
	factory := grpc.NewGrpcConnectionFactory(url).WithInsecure()
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		factory = factory.WithPerRPCCredentials(auth.NewOauthAccessWithoutTLS(token))
	}
	return factory.WithRetry().WithBlock().Build(ctx)
}

func CreateGatewayConnecton(ctx context.Context, url string, token *oauth2.Token) (*ggrpc.ClientConn, error) {
	factory := grpc.NewGrpcConnectionFactory(url).WithTransportCredentials()
	if token != nil {
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		factory = factory.WithPerRPCCredentials(oauth.NewOauthAccess(token))
	}
	return factory.WithRetry().WithBlock().Build(ctx)
}
