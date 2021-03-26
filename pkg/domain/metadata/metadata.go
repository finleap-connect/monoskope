package domain

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/grpc/metadata"
)

const (
	componentInformationKey = "component_information"
	rolebindingsKey         = "user_role_bindings"
)

// ComponentInformation are information about a service/component.
type ComponentInformation struct {
	Name    string
	Version string
	Commit  string
}

// UserInformation are identifying information about a user.
type UserInformation struct {
	Id     uuid.UUID
	Name   string
	Email  string
	Issuer string
}

// domainMetadataManager is a domain specific metadata manager.
type DomainMetadataManager struct {
	es.MetadataManager
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (*DomainMetadataManager, error) {
	m := &DomainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
	}

	// Get the grpc metadata from incoming context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		data := make(map[string]string)
		for k, v := range md {
			data[k] = v[0] // typically only the first and only value of that is relevant
		}
		m.SetMetadata(data)
	}

	if _, exists := m.Get(componentInformationKey); !exists {
		if err := m.SetComponentInformation(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// SetComponentInformation sets the ComponentInformation about the currently executing service/component.
func (m *DomainMetadataManager) SetComponentInformation() error {
	return m.SetObject(componentInformationKey, &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
}

func (m *DomainMetadataManager) SetRoleBindings(roleBindings []*projections.UserRoleBinding) error {
	return m.SetObject(rolebindingsKey, roleBindings)
}

func (m *DomainMetadataManager) GetRoleBindings() ([]*projections.UserRoleBinding, error) {
	roleBindings := make([]*projections.UserRoleBinding, 0)
	err := m.GetObject(rolebindingsKey, &roleBindings)
	return roleBindings, err
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *DomainMetadataManager) SetUserInformation(userInformation *UserInformation) {
	m.Set(gateway.HeaderAuthName, userInformation.Name)
	m.Set(gateway.HeaderAuthEmail, userInformation.Email)
	m.Set(gateway.HeaderAuthIssuer, userInformation.Issuer)
	m.Set(gateway.HeaderAuthId, userInformation.Id.String())
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *DomainMetadataManager) GetUserInformation() *UserInformation {
	userInfo := &UserInformation{}
	if header, ok := m.Get(gateway.HeaderAuthName); ok {
		userInfo.Name = header
	}
	if header, ok := m.Get(gateway.HeaderAuthEmail); ok {
		userInfo.Email = header
	}
	if header, ok := m.Get(gateway.HeaderAuthIssuer); ok {
		userInfo.Issuer = header
	}
	if header, ok := m.Get(gateway.HeaderAuthId); ok {
		id, err := uuid.Parse(header)
		if err == nil {
			userInfo.Id = id
		}
	}
	return userInfo
}

// GetOutgoingGrpcContext returns a new context enriched with the metadata of this manager.
func (m *DomainMetadataManager) GetOutgoingGrpcContext() context.Context {
	return metadata.NewOutgoingContext(m.MetadataManager.GetContext(), metadata.New(m.GetMetadata()))
}
