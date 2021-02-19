package domain

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/grpc/metadata"
)

const (
	userInformationKey      = "user_information"
	componentInformationKey = "component_information"
)

// ComponentInformation are information about a service/component.
type ComponentInformation struct {
	Name    string
	Version string
	Commit  string
}

// UserInformation are identifying information about a user.
type UserInformation struct {
	Email   string
	Subject string
	Issuer  string
}

// domainMetadataManager is a domain specific metadata manager.
type domainMetadataManager struct {
	es.MetadataManager
}

type DomainMetadataManager interface {
	es.MetadataManager
	SetComponentInformation() error
	SetUserInformation(userInformation *UserInformation) error
	GetUserInformation() (*UserInformation, error)
	GetOutgoingGrpcContext() context.Context
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (DomainMetadataManager, error) {
	m := &domainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		data := make(map[string]string)
		for k, v := range md {
			data[k] = v[0]
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
func (m *domainMetadataManager) SetComponentInformation() error {
	return m.SetObject(componentInformationKey, &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *domainMetadataManager) SetUserInformation(userInformation *UserInformation) error {
	return m.SetObject(userInformationKey, userInformation)
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *domainMetadataManager) GetUserInformation() (*UserInformation, error) {
	userInfo := &UserInformation{}
	err := m.GetObject(userInformationKey, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, err
}

// GetOutgoingGrpcContext returns a new context enriched with the metadata of this manager.
func (m *domainMetadataManager) GetOutgoingGrpcContext() context.Context {
	return metadata.NewOutgoingContext(m.MetadataManager.GetContext(), metadata.New(m.GetMetadata()))
}
