package errors

import "errors"

// Domain Errors
var (
	// ErrUnauthorized is when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("no events to append")
)
