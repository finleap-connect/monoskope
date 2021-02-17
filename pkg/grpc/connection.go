package grpc

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type grpcConnectionFactory struct {
	opts []grpc.DialOption
	url  string
}

func NewGrpcConnectionFactory(url string) grpcConnectionFactory {
	return grpcConnectionFactory{
		url: url,
	}
}

func NewGrpcConnectionFactoryWithDefaults(url string) grpcConnectionFactory {
	return NewGrpcConnectionFactory(url).
		WithInsecure().
		WithRetry().
		WithBlock()
}

func (factory grpcConnectionFactory) WithInsecure() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithInsecure())
	return factory
}

func (factory grpcConnectionFactory) WithPerRPCCredentials(creds credentials.PerRPCCredentials) grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithPerRPCCredentials(creds))
	return factory
}

func (factory grpcConnectionFactory) WithTransportCredentials() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	return factory
}

func (factory grpcConnectionFactory) WithBlock() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}
	factory.opts = append(factory.opts, grpc.WithBlock())
	return factory
}

func (factory grpcConnectionFactory) WithRetry() grpcConnectionFactory {
	if factory.opts == nil {
		factory.opts = make([]grpc.DialOption, 0)
	}

	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(10 * time.Millisecond)),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...),
		grpc_retry.WithMax(5),
	}

	factory.opts = append(factory.opts, grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)))
	factory.opts = append(factory.opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))

	return factory
}

func (factory grpcConnectionFactory) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, factory.url, factory.opts...)
}

func (factory grpcConnectionFactory) ConnectWithTimeout(ctx context.Context, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return factory.Connect(ctx)
}
