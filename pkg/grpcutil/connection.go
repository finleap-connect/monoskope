package grpcutil

import (
	"context"

	"google.golang.org/grpc"
)

// Creates a gRPC connection with insecure.
// opts can be passed as nil.
func CreateInsecureGrpcConnecton(ctx context.Context, url string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	if opts == nil {
		opts = make([]grpc.DialOption, 0)
	}
	opts = append(opts, grpc.WithInsecure())
	return CreateGrpcConnection(ctx, url, opts)
}

// Creates a gRPC connection.
// opts can be passed as nil.
func CreateGrpcConnection(ctx context.Context, url string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	if opts == nil {
		opts = make([]grpc.DialOption, 0)
	}
	opts = append(opts, grpc.WithBlock())
	return grpc.DialContext(ctx, url, opts...)
}
