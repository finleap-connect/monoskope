package grpc

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrInvalidToken is the gRPC error response for invalid auth token
	ErrInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

func ErrInternal(msg string) error {
	return status.Error(codes.Internal, msg)
}

func ErrMandatory(fields ...string) error {
	var msg string
	if len(fields) == 1 {
		msg = fmt.Sprintf("%s is mandatory", fields[0])
	} else {
		msg = fmt.Sprintf("%s are mandatory", strings.Join(fields, ", "))
	}
	return status.Error(codes.InvalidArgument, msg)
}

func ErrMalformed(fields ...string) error {
	var msg string
	if len(fields) == 1 {
		msg = fmt.Sprintf("%s is malformed", fields[0])
	} else {
		msg = fmt.Sprintf("%s are malformed", strings.Join(fields, ", "))
	}
	return status.Error(codes.InvalidArgument, msg)
}

func ErrInvalidArgument(err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
}

func IsStatus(code codes.Code, err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return s.Code() == code
}
