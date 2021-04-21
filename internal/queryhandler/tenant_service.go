package queryhandler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	grpcUtil "gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"google.golang.org/grpc"
)

// tenantServiceServer is the implementation of the TenantService API
type tenantServiceServer struct {
	api.UnimplementedTenantServer

	repo repositories.ReadOnlyTenantRepository
}

// NewTenantServiceServer returns a new configured instance of tenantServiceServer
func NewTenantServer(tenantRepo repositories.ReadOnlyTenantRepository) *tenantServiceServer {
	return &tenantServiceServer{
		repo: tenantRepo,
	}
}

func NewTenantClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.TenantClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewTenantClient(conn), nil
}

// GetById returns the tenant found by the given id.
func (s *tenantServiceServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.Tenant, error) {
	tenant, err := s.repo.ByTenantId(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}
	return tenant.Proto(), nil
}

// GetByName returns the tenant found by the given name.
func (s *tenantServiceServer) GetByName(ctx context.Context, name *wrappers.StringValue) (*projections.Tenant, error) {
	tenant, err := s.repo.ByName(ctx, name.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return tenant.Proto(), nil
}

func (s *tenantServiceServer) GetAll(request *api.GetAllRequest, stream api.Tenant_GetAllServer) error {
	users, err := s.repo.GetAll(stream.Context(), request.GetIncludeDeleted())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, user := range users {
		err := stream.Send(user.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}
