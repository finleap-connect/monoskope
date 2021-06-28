package repositories

import (
	"context"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type remoteClusterRepository struct {
	clusterClient api.ClusterClient
}

// NewRemoteClusterRepository creates a repository for reading user projections.
func NewRemoteClusterRepository(clusterClient api.ClusterClient) ReadOnlyClusterRepository {
	return &remoteClusterRepository{
		clusterClient: clusterClient,
	}
}

func (r *remoteClusterRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Cluster, error) {
	panic("not implemented")
}

func (r *remoteClusterRepository) GetBootstrapToken(ctx context.Context, id string) (string, error) {
	panic("not implemented")
}

func (r *remoteClusterRepository) ByClusterId(ctx context.Context, id string) (*projections.Cluster, error) {
	clusterProto, err := r.clusterClient.GetById(ctx, wrapperspb.String(id))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}
	return &projections.Cluster{Cluster: clusterProto}, nil
}

func (r *remoteClusterRepository) ByClusterName(ctx context.Context, name string) (*projections.Cluster, error) {
	clusterProto, err := r.clusterClient.GetByName(ctx, wrapperspb.String(name))
	if err != nil {
		return nil, errors.TranslateFromGrpcError(err)
	}
	return &projections.Cluster{Cluster: clusterProto}, nil
}
