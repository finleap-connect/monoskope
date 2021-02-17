package domain

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

const (
	userInformationKey      = "userInformationKey"
	componentInformationKey = "componentInformationKey"
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
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (DomainMetadataManager, error) {
	m := &domainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
	}
	if err := m.SetComponentInformation(); err != nil {
		return nil, err
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
