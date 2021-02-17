package errors

import "errors"

var (
	// ErrUnauthorized is when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrUserNotFound is when a user is not known to the system.
	ErrUserNotFound = errors.New("user not found")
)
