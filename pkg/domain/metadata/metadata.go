package domain

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	userInformationKey      evs.EventMetadataKey
	componentInformationKey evs.EventMetadataKey
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
	evs.MetadataManager
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) *domainMetadataManager {
	m := &domainMetadataManager{
		evs.NewMetadataManagerFromContext(ctx),
	}
	return m.SetComponentInformation()
}

// SetComponentInformation sets the ComponentInformation about the currently executing service/component.
func (m *domainMetadataManager) SetComponentInformation() *domainMetadataManager {
	m.Set(evs.EventMetadataKey(componentInformationKey), &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
	return m
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *domainMetadataManager) SetUserInformation(userInformation *UserInformation) *domainMetadataManager {
	m.Set(userInformationKey, userInformation)
	return m
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *domainMetadataManager) GetUserInformation() (*UserInformation, error) {
	iface, ok := m.Get(userInformationKey)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	userId, ok := iface.(UserInformation)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	return &userId, nil
}
