package grpcutil

import (
	"context"

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

func NewGrpcConnectionFactoryWithOpts(url string, opts []grpc.DialOption) grpcConnectionFactory {
	return grpcConnectionFactory{
		url:  url,
		opts: opts,
	}
}

func (factory grpcConnectionFactory) Url(url string) grpcConnectionFactory {
	factory.url = url
	return factory
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

func (factory grpcConnectionFactory) Build(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, factory.url, factory.opts...)
}
