// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"

	es_errors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrNotFound is returned when an aggregate is not known to the system.
	ErrNotFound = errors.New("aggregate not found")
	// ErrDeleted is returned when an aggregate has been deleted.
	ErrDeleted = errors.New("aggregate is deleted")
	// ErrUnknownAggregateType is returned returned when an unknown or invalid aggregate type is used.
	ErrUnknownAggregateType = errors.New("aggregate type unknown")

	// ErrUnauthorized is returned when a requested command can not be executed because the current user is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrUserNotFound is returned when a user is not known to the system.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists is returned when a user does already exist.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserRoleBindingAlreadyExists is returned when a userrolebinding does already exist.
	ErrUserRoleBindingAlreadyExists = errors.New("userrolebinding already exists")
	// ErrUserRoleBindingNotFound is returned when a binding by that name or id cannot be found.
	ErrUserRoleBindingNotFound = errors.New("userrolebinding not found")

	// ErrTenantNotFound is returned when a tenant is not known to the system.
	ErrTenantNotFound = errors.New("tenant not found")
	// ErrTenantAlreadyExists is returned when a tenant does already exist.
	ErrTenantAlreadyExists = errors.New("tenant already exists")

	// ErrClusterRegistrationNotFound is returned when a cluster registration is not known to the system.
	ErrClusterRegistrationNotFound = errors.New("cluster registration not found")
	// ErrClusterRegistrationAlreadyExists is returned when an aggregate does already exist.
	ErrClusterRegistrationAlreadyExists = errors.New("cluster registration already exists")

	// ErrClusterNotFound is returned when a cluster is not known to the system.
	ErrClusterNotFound = errors.New("cluster not found")
	// ErrClusterAlreadyExists is returned when a cluster does already exist.
	ErrClusterAlreadyExists = errors.New("cluster already exists")

	// ErrCertificateAlreadyExists is returned when a cert does already exist.
	ErrCertificateAlreadyExists = errors.New("certificate already exists")

	// ErrTenantClusterBindingAlreadyExists is returned when a tenant-cluster-binding does already exist.
	ErrTenantClusterBindingAlreadyExists = errors.New("tenant already has access to that cluster")
	// ErrTenantClusterBindingNotFound is returned when a tenant-cluster-binding could not be found.
	ErrTenantClusterBindingNotFound = errors.New("no cluster access found for the given cluster and tenant")
)

var (
	errorMap = map[codes.Code][]error{
		codes.NotFound: {
			ErrUserNotFound,
			ErrTenantNotFound,
			ErrClusterRegistrationNotFound,
			ErrClusterNotFound,
			es_errors.ErrProjectionNotFound,
		},
		codes.AlreadyExists: {
			ErrUserAlreadyExists,
			ErrTenantAlreadyExists,
			ErrClusterRegistrationAlreadyExists,
			ErrClusterAlreadyExists,
			ErrCertificateAlreadyExists,
			ErrTenantClusterBindingAlreadyExists,
		},
		codes.PermissionDenied: {ErrUnauthorized},
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

// Returns a gRPC status error with InvalidArgument code
func ErrInvalidArgument(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// Returns a gRPC status error with FailedPrecondition code
func ErrFailedPrecondition(msg string) error {
	return status.Error(codes.FailedPrecondition, msg)
}
