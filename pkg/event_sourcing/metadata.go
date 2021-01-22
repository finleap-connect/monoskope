package event_sourcing

import (
	"context"
)

type EventMetadataKey int

var metadataKey EventMetadataKey

type MetadataManager interface {
	Get(EventMetadataKey) (interface{}, bool)
	Set(EventMetadataKey, interface{}) MetadataManager
	GetContext() context.Context
}

type metadataManager struct {
	ctx  context.Context
	data map[EventMetadataKey]interface{}
}

func NewMetadataManager(ctx context.Context) MetadataManager {
	b := &metadataManager{
		ctx:  ctx,
		data: make(map[EventMetadataKey]interface{}),
	}

	d, ok := ctx.Value(metadataKey).(map[EventMetadataKey]interface{})
	if ok {
		b.data = d
	}

	return b
}

func (b *metadataManager) Get(key EventMetadataKey) (interface{}, bool) {
	v, ok := b.data[key]
	return v, ok
}

func (b *metadataManager) Set(key EventMetadataKey, value interface{}) MetadataManager {
	b.data[key] = value
	return b
}

func (b *metadataManager) GetContext() context.Context {
	return context.WithValue(b.ctx, metadataKey, b.data)
}
