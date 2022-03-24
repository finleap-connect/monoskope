package middleware

import (
	"google.golang.org/grpc"
)

type GRPCMiddleware interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
}
