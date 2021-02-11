package errors

import "errors"

// Authorization
var (
	// ErrUnauthorized is when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")
)
