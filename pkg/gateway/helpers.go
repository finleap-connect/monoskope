package gateway

import (
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateGatewayAuthedConnecton(url string, transportCredentials credentials.TransportCredentials, token *oauth2.Token) (*grpc.ClientConn, error) {
	perRPC := oauth.NewOauthAccess(token)

	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.NewOauthAccess requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(transportCredentials),
	}

	opts = append(opts, grpc.WithBlock())
	return grpc.Dial(url, opts...)
}

func CreateGatewayConnecton(url string, transportCredentials credentials.TransportCredentials) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCredentials),
	}

	opts = append(opts, grpc.WithBlock())
	return grpc.Dial(url, opts...)
}
