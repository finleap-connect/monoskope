package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrUnauthorized is when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrUserNotFound is when a user is not known to the system.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is when a user does already exist.
	ErrUserAlreadyExists = errors.New("user already exists")
)

var (
	errorMap = map[codes.Code]error{
		codes.NotFound:         ErrUserNotFound,
		codes.PermissionDenied: ErrUnauthorized,
		codes.AlreadyExists:    ErrUserAlreadyExists,
	}
	reverseErrorMap = reverseMap(errorMap)
)

func reverseMap(m map[codes.Code]error) map[error]codes.Code {
	n := make(map[error]codes.Code)
	for k, v := range m {
		n[v] = k
	}
	return n
}

// TranslateFromGrpcError takes an error that could be an error received
// via gRPC and converts it to a specific domain error to make it easily comparable.
func TranslateFromGrpcError(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	if mappedErr, ok := errorMap[s.Code()]; ok {
		return mappedErr
	}

	return err
}

// TranslateToGrpcError converts an error to a gRPC status error.
func TranslateToGrpcError(err error) error {
	if code, ok := reverseErrorMap[err]; ok {
		return status.Error(code, err.Error())
	}
	return err
}
