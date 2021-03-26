package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrNotFound is when an aggregate is not known to the system.
	ErrNotFound = errors.New("not found")
	// ErrDeleted is when an aggregate has been deleted.
	ErrDeleted = errors.New("deleted")

	// ErrUnauthorized is when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrUserNotFound is when a user is not known to the system.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists is when a user does already exist.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserRoleBindingAlreadyExists is when a userrolebinding does already exist.
	ErrUserRoleBindingAlreadyExists = errors.New("userrolebinding already exists")

	// ErrTenantNotFound is when a tenant is not known to the system.
	ErrTenantNotFound = errors.New("tenant not found")
	// ErrTenantAlreadyExists is when a tenant does already exist.
	ErrTenantAlreadyExists = errors.New("tenant already exists")
)

var (
	errorMap = map[codes.Code][]error{
		codes.NotFound:         {ErrUserNotFound, ErrTenantNotFound},
		codes.PermissionDenied: {ErrUnauthorized},
		codes.AlreadyExists:    {ErrUserAlreadyExists, ErrTenantAlreadyExists},
	}
	reverseErrorMap = reverseMap(errorMap)
)

func reverseMap(m map[codes.Code][]error) map[error]codes.Code {
	n := make(map[error]codes.Code)
	for k, v := range m {
		for _, e := range v {
			n[e] = k
		}
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
		for _, e := range mappedErr {
			if TranslateToGrpcError(e).Error() == err.Error() {
				return e
			}
		}
	}

	return err
}

// TranslateToGrpcError converts an error to a gRPC status error.
func TranslateToGrpcError(err error) error {
	if code, ok := reverseErrorMap[err]; ok {
		return status.Error(code, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}

func ErrInvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}
