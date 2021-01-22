package domain

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

const (
	UserId                  evs.EventMetadataKey = "UserId"
	ComponentInformationKey evs.EventMetadataKey = "ComponentInformation"
)

type ComponentInformation struct {
	Name    string
	Version string
	Commit  string
}

type metadataBuilder struct {
	ctx  context.Context
	data map[evs.EventMetadataKey]interface{}
}

func NewMetadataBuilder(ctx context.Context) *metadataBuilder {
	b := &metadataBuilder{
		ctx:  ctx,
		data: make(map[evs.EventMetadataKey]interface{}),
	}
	return b.SetComponentInformation()
}

func (b *metadataBuilder) SetComponentInformation() *metadataBuilder {
	b.data[ComponentInformationKey] = &ComponentInformation{
		Name:    version.Name,
		Version: version.Version,
		Commit:  version.Commit,
	}
	return b
}

func (b *metadataBuilder) SetUserId(userId uuid.UUID) *metadataBuilder {
	b.data[UserId] = userId
	return b
}

func (b *metadataBuilder) Apply() context.Context {
	return context.WithValue(b.ctx, evs.EventMedataData, b.data)
}
