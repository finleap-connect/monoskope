package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"
)

func CreateInsecureGatewayConnecton(ctx context.Context, url string) (*ggrpc.ClientConn, error) {
	factory := grpc.NewGrpcConnectionFactory(url).WithInsecure()
	return factory.WithRetry().WithBlock().Connect(ctx)
}
