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

package repositories

import (
	"context"
	"fmt"
	"io"

	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type remoteUserRepository struct {
	userService api.UserClient
}

// NewRemoteUserRepository creates a repository for reading user projections.
func NewRemoteUserRepository(userService api.UserClient) ReadOnlyUserRepository {
	return &remoteUserRepository{
		userService: userService,
	}
}

func (r *remoteUserRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.User, error) {
	return nil, fmt.Errorf("not implemented")
}

// ById searches for the a user projection by it's id.
func (r *remoteUserRepository) ByUserId(ctx context.Context, id uuid.UUID) (*projections.User, error) {
	userProto, err := r.userService.GetById(ctx, wrapperspb.String(id.String()))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	user := &projections.User{User: userProto}

	// Find roles of user
	stream, err := r.userService.GetRoleBindingsById(ctx, wrapperspb.String(user.Id))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	for {
		// Read next event
		proto, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, errors.TranslateFromGrpcError(err)
		}

		user.Roles = append(user.Roles, proto)
	}

	return user, nil
}

// ByEmail searches for the a user projection by it's email address.
func (r *remoteUserRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
	userProto, err := r.userService.GetByEmail(ctx, wrapperspb.String(email))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	user := &projections.User{User: userProto}

	// Find roles of user
	stream, err := r.userService.GetRoleBindingsById(ctx, wrapperspb.String(user.Id))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	for {
		// Read next event
		proto, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return nil, errors.TranslateFromGrpcError(err)
		}

		user.Roles = append(user.Roles, proto)
	}

	return user, nil
}

func (r *remoteUserRepository) GetCount(ctx context.Context, includeDeleted bool) (int, error) {
	return -1, fmt.Errorf("not implemented")
}
