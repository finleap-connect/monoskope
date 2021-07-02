package repositories

import (
	"context"
	"io"

	"github.com/google/uuid"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
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
	panic("not implemented")
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
