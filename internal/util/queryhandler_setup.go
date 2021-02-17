package util

import (
	"context"
	"time"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

func NewUserServiceClient(queryHandlerAddr string) (*grpc.ClientConn, api.UserServiceClient, error) {
	// Create EventStore client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpcUtil.
		NewGrpcConnectionFactory(queryHandlerAddr).
		WithInsecure().
		WithRetry().
		WithBlock().
		Build(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, api.NewUserServiceClient(conn), nil
}
