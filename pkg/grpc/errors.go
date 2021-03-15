package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrInvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}
