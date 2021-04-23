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
	"google.golang.org/protobuf/types/known/emptypb"
)

// clusterregistrationServer is the implementation of the ClusterRegistrationService API
type clusterregistrationServer struct {
	api.UnimplementedClusterRegistrationServer

	repo repositories.ReadOnlyClusterRegistrationRepository
}

// NewClusterRegistrationServiceServer returns a new configured instance of clusterregistrationServiceServer
func NewClusterRegistrationServer(clusterregistrationRepo repositories.ReadOnlyClusterRegistrationRepository) *clusterregistrationServer {
	return &clusterregistrationServer{
		repo: clusterregistrationRepo,
	}
}

func NewClusterRegistrationClient(ctx context.Context, queryHandlerAddr string) (*grpc.ClientConn, api.ClusterRegistrationClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(queryHandlerAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, errors.TranslateToGrpcError(err)
	}

	return conn, api.NewClusterRegistrationClient(conn), nil
}

func (s *clusterregistrationServer) GetPending(empty *emptypb.Empty, stream api.ClusterRegistration_GetPendingServer) error {
	projections, err := s.repo.GetPending(stream.Context())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, projection := range projections {
		err := stream.Send(projection.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}

// GetById returns the clusterregistration found by the given id.
func (s *clusterregistrationServer) GetById(ctx context.Context, id *wrappers.StringValue) (*projections.ClusterRegistration, error) {
	projection, err := s.repo.ByClusterRegistrationId(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}
	return projection.Proto(), nil
}

// GetByName returns the clusterregistration found by the given name.
func (s *clusterregistrationServer) GetByName(ctx context.Context, name *wrappers.StringValue) (*projections.ClusterRegistration, error) {
	projection, err := s.repo.ByName(ctx, name.GetValue())
	if err != nil {
		return nil, errors.TranslateToGrpcError(err)
	}
	return projection.Proto(), nil
}

func (s *clusterregistrationServer) GetAll(request *api.GetAllRequest, stream api.ClusterRegistration_GetAllServer) error {
	projections, err := s.repo.GetAll(stream.Context(), request.GetIncludeDeleted())
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	for _, projection := range projections {
		err := stream.Send(projection.Proto())
		if err != nil {
			return errors.TranslateToGrpcError(err)
		}
	}
	return nil
}
