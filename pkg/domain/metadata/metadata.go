package domain

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
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
	Name   string
	Email  string
	Issuer string
}

// domainMetadataManager is a domain specific metadata manager.
type domainMetadataManager struct {
	es.MetadataManager
}

type DomainMetadataManager interface {
	es.MetadataManager
	SetComponentInformation() error
	SetUserInformation(userInformation *UserInformation)
	GetUserInformation() *UserInformation
	GetOutgoingGrpcContext() context.Context
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (DomainMetadataManager, error) {
	m := &domainMetadataManager{
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
func (m *domainMetadataManager) SetComponentInformation() error {
	return m.SetObject(componentInformationKey, &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *domainMetadataManager) SetUserInformation(userInformation *UserInformation) {
	m.Set(gateway.HeaderAuthName, userInformation.Name)
	m.Set(gateway.HeaderAuthEmail, userInformation.Email)
	m.Set(gateway.HeaderAuthIssuer, userInformation.Issuer)
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *domainMetadataManager) GetUserInformation() *UserInformation {
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
	return userInfo
}

// GetOutgoingGrpcContext returns a new context enriched with the metadata of this manager.
func (m *domainMetadataManager) GetOutgoingGrpcContext() context.Context {
	return metadata.NewOutgoingContext(m.MetadataManager.GetContext(), metadata.New(m.GetMetadata()))
}
