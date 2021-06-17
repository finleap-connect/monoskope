package repositories

import (
	"context"

	"github.com/google/uuid"
	apiProjections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type clusterRepository struct {
	*domainRepository
}

// ClusterRepository is a repository for reading and writing cluster projections.
type ClusterRepository interface {
	es.Repository
	ReadOnlyClusterRepository
	WriteOnlyClusterRepository
}

// ReadOnlyClusterRepository is a repository for reading cluster projections.
type ReadOnlyClusterRepository interface {
	// ById searches for the a tenant projection by it's id.
	ByClusterId(context.Context, string) (*projections.Cluster, error)
	// ByName searches for the a tenant projection by it's name
	ByClusterName(context.Context, string) (*projections.Cluster, error)
	// GetAll searches for all known clusters.
	GetAll(context.Context, bool) ([]*projections.Cluster, error)
	// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID
	GetBootstrapToken(context.Context, string) (string, error)
	// GetCertificate returns the certificate issued for the m8 operator of the cluster with the given UUID
	GetCertificate(context.Context, string) (*apiProjections.Certificate, error)
}

// WriteOnlyClusterRepository is a repository for writing cluster projections.
type WriteOnlyClusterRepository interface {
}

// NewClusterRepository creates a repository for reading and writing cluster projections.
func NewClusterRepository(repository es.Repository, userRepo UserRepository) ClusterRepository {
	return &clusterRepository{
		domainRepository: NewDomainRepository(repository, userRepo),
	}
}

// ById searches for a cluster by its id.
func (r *clusterRepository) ByClusterId(ctx context.Context, id string) (*projections.Cluster, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	cluster, ok := projection.(*projections.Cluster)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	err = r.addMetadata(ctx, cluster.DomainProjection)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// ByClusterName searches for a cluster projection by its name.
func (r *clusterRepository) ByClusterName(ctx context.Context, clusterName string) (*projections.Cluster, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, c := range ps {
		if clusterName == c.Name {
			return c, nil
		}
	}

	return nil, errors.ErrClusterNotFound
}

// GetAll searches for all cluster projections.
func (r *clusterRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Cluster, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	var clusters []*projections.Cluster
	for _, p := range ps {
		if c, ok := p.(*projections.Cluster); ok {
			if !c.GetDeleted().IsValid() || includeDeleted {
				// Add metadata
				err = r.addMetadata(ctx, c.DomainProjection)
				if err != nil {
					return nil, err
				}
				clusters = append(clusters, c)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}
	return clusters, nil
}

// GetBootstrapToken returns the bootstrap token for a cluster with the given UUID.
func (r *clusterRepository) GetBootstrapToken(ctx context.Context, id string) (string, error) {
	cluster, err := r.ByClusterId(ctx, id)
	if err != nil {
		return "", err
	}
	return cluster.BootstrapToken, nil
}

// GetClusterCertificates returns the m8 CA and the certificate issued for the m8 operator of the cluster with the given UUID
func (r *clusterRepository) GetCertificate(ctx context.Context, id string) (*apiProjections.Certificate, error) {
	cluster, err := r.ByClusterId(ctx, id)
	if err != nil {
		return nil, err
	}
	return cluster.GetCertificate(), nil
}
