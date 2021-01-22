package domain

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var (
	userMailKey             evs.EventMetadataKey
	componentInformationKey evs.EventMetadataKey
)

type ComponentInformation struct {
	Name    string
	Version string
	Commit  string
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

func (b *domainMetadataManager) SetUserEmail(userMail string) *domainMetadataManager {
	b.metadataManager.Set(userMailKey, userMail)
	return b
}

func (b *domainMetadataManager) GetUserEmail() (string, error) {
	iface, ok := b.metadataManager.Get(userMailKey)
	if !ok {
		return "", fmt.Errorf("not found")
	}

	userId, ok := iface.(string)
	if !ok {
		return "", fmt.Errorf("invalid type")
	}
	return userId, nil
}

func (b *domainMetadataManager) GetContext() context.Context {
	return b.metadataManager.GetContext()
}
