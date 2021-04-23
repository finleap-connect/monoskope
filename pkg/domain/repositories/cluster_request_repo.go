package repositories

import (
	"context"

	"github.com/google/uuid"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type clusterregistrationRepository struct {
	*approvableRepository
}

// ClusterRegistrationRepository is a repository for reading and writing clusterregistration projections.
type ClusterRegistrationRepository interface {
	es.Repository
	ReadOnlyClusterRegistrationRepository
	WriteOnlyClusterRegistrationRepository
}

// ReadOnlyClusterRegistrationRepository is a repository for reading clusterregistration projections.
type ReadOnlyClusterRegistrationRepository interface {
	// ById searches for the a clusterregistration projection by it's id.
	ByClusterRegistrationId(context.Context, string) (*projections.ClusterRegistration, error)
	// ByName searches for the a clusterregistration projection by it's name
	ByName(context.Context, string) (*projections.ClusterRegistration, error)
	// GetAll searches for all clusterregistration projections.
	GetAll(context.Context, bool) ([]*projections.ClusterRegistration, error)
	// GetPending searches for the pending clusterregistration projections.
	GetPending(context.Context) ([]*projections.ClusterRegistration, error)
}

// WriteOnlyClusterRegistrationRepository is a repository for writing clusterregistration projections.
type WriteOnlyClusterRegistrationRepository interface {
}

// NewClusterRegistrationRepository creates a repository for reading and writing clusterregistration projections.
func NewClusterRegistrationRepository(repository es.Repository, userRepo UserRepository) ClusterRegistrationRepository {
	return &clusterregistrationRepository{
		approvableRepository: NewApprovableRepository(repository, userRepo),
	}
}

// ByClusterRegistrationId searches for a clusterregistration projection by its id.
func (r *clusterregistrationRepository) ByClusterRegistrationId(ctx context.Context, id string) (*projections.ClusterRegistration, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	clusterregistration, ok := projection.(*projections.ClusterRegistration)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	err = r.addMetadata(ctx, clusterregistration.ApprovableProjection)
	if err != nil {
		return nil, err
	}
	return clusterregistration, nil
}

// ByClusterRegistrationName searches for a clusterregistration projection by its name.
func (r *clusterregistrationRepository) ByName(ctx context.Context, name string) (*projections.ClusterRegistration, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, t := range ps {
		if name == t.Name {
			return t, nil
		}
	}

	return nil, errors.ErrClusterRegistrationNotFound
}

// All searches for the a clusterregistration projections.
func (r *clusterregistrationRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.ClusterRegistration, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var clusterregistrations []*projections.ClusterRegistration
	for _, p := range ps {
		if t, ok := p.(*projections.ClusterRegistration); ok {
			if !t.GetDeleted().IsValid() || includeDeleted {
				// Add metadata
				err = r.addMetadata(ctx, t.ApprovableProjection)
				if err != nil {
					return nil, err
				}
				clusterregistrations = append(clusterregistrations, t)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}
	return clusterregistrations, nil
}

// GetPending searches for the pending clusterregistration projections.
func (r *clusterregistrationRepository) GetPending(ctx context.Context) ([]*projections.ClusterRegistration, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	var clusterregistrations []*projections.ClusterRegistration
	for _, t := range ps {
		if t.Status == projectionsApi.ClusterRegistration_REQUESTED {
			clusterregistrations = append(clusterregistrations, t)
		}
	}
	return clusterregistrations, nil
}
