package gateway

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	ggrpc "google.golang.org/grpc"
)

func CreateInsecureConnection(ctx context.Context, url string) (*ggrpc.ClientConn, error) {
	return grpc.NewGrpcConnectionFactory(url).
		WithInsecure().
		WithRetry().
		WithBlock().
		Connect(ctx)
}
