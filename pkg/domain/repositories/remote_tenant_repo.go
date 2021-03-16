package repositories

import (
	"context"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type remoteTenantRepository struct {
	tenantService api.TenantServiceClient
}

// NewRemoteTenantRepository creates a repository for reading tenant projections.
func NewRemoteTenantRepository(tenantService api.TenantServiceClient) ReadOnlyTenantRepository {
	return &remoteTenantRepository{
		tenantService: tenantService,
	}
}

// ById searches for the a tenant projection by its id.
func (r *remoteTenantRepository) ByTenantId(ctx context.Context, id string) (*projections.Tenant, error) {
	tenantProto, err := r.tenantService.GetById(ctx, wrapperspb.String(id))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	tenant := &projections.Tenant{Tenant: tenantProto}

	// add users to tenant

	// for {
	// 	// Read next event
	// 	proto, err := stream.Recv()

	// 	// End of stream
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil { // Some other error
	// 		return nil, errors.TranslateFromGrpcError(err)
	// 	}

	// 	user.Roles = append(user.Roles, proto)
	// }

	return tenant, nil
}

// ByEmail searches for the a user projection by it's email address.
func (r *remoteTenantRepository) ByName(ctx context.Context, name string) (*projections.Tenant, error) {
	tenantProto, err := r.tenantService.GetByName(ctx, wrapperspb.String(name))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}

	tenant := &projections.Tenant{Tenant: tenantProto}

	// add users to tenant
	// stream, err := r.userService.GetRoleBindingsById(ctx, wrapperspb.String(user.Id))
	// if err != nil {
	// 	return nil, errors.TranslateFromGrpcError(err)
	// }

	// for {
	// 	// Read next event
	// 	proto, err := stream.Recv()

	// 	// End of stream
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil { // Some other error
	// 		return nil, errors.TranslateFromGrpcError(err)
	// 	}

	// 	user.Roles = append(user.Roles, proto)
	// }

	return tenant, nil
}
