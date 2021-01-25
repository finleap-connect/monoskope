package domain

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var (
	userInformationKey      evs.EventMetadataKey
	componentInformationKey evs.EventMetadataKey
)

type ComponentInformation struct {
	Name    string
	Version string
	Commit  string
}

type UserInformation struct {
	Email   string
	Subject string
	Issuer  string
}

type domainMetadataManager struct {
	metadataManager evs.MetadataManager
}

func NewDomainMetadataManager(ctx context.Context) *domainMetadataManager {
	b := &domainMetadataManager{
		metadataManager: evs.NewMetadataManager(ctx),
	}
	return b.SetComponentInformation()
}

func (b *domainMetadataManager) SetComponentInformation() *domainMetadataManager {
	b.metadataManager.Set(evs.EventMetadataKey(componentInformationKey), &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	})
	return b
}

func (b *domainMetadataManager) SetUserInformation(userInformation *UserInformation) *domainMetadataManager {
	b.metadataManager.Set(userInformationKey, userInformation)
	return b
}

func (b *domainMetadataManager) GetUserInformation() (*UserInformation, error) {
	iface, ok := b.metadataManager.Get(userInformationKey)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	userId, ok := iface.(UserInformation)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	return &userId, nil
}

func (b *domainMetadataManager) GetContext() context.Context {
	return b.metadataManager.GetContext()
}
