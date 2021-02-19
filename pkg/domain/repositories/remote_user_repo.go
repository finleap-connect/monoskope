package repositories

import (
	"context"
	"io"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type remoteUserRepository struct {
	userService api.UserServiceClient
}

// NewRemoteUserRepository creates a repository for reading user projections.
func NewRemoteUserRepository(userService api.UserServiceClient) ReadOnlyUserRepository {
	return &remoteUserRepository{
		userService: userService,
	}
}

// ById searches for the a user projection by it's id.
func (r *remoteUserRepository) ByUserId(ctx context.Context, id string) (*projections.User, error) {
	userProto, err := r.userService.GetById(ctx, wrapperspb.String(id))
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
