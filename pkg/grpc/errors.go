package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrInternal(msg string) error {
	return status.Error(codes.Internal, msg)
}

func ErrInvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

func IsStatus(code codes.Code, err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return s.Code() == code
}
