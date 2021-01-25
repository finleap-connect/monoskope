package event_sourcing

import (
	"context"
	"fmt"
)

type EventMetadataKey int

var metadataKey EventMetadataKey

type MetadataManager interface {
	GetContext() context.Context

	Get(EventMetadataKey) (interface{}, bool)
	GetBool(EventMetadataKey) (bool, error)
	GetString(EventMetadataKey) (string, error)

	Set(EventMetadataKey, interface{}) MetadataManager
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

func (b *metadataManager) GetString(key EventMetadataKey) (string, error) {
	iface, ok := b.Get(key)
	if !ok {
		return "", fmt.Errorf("not found")
	}

	stringValue, ok := iface.(string)
	if !ok {
		return "", fmt.Errorf("invalid type")
	}
	return stringValue, nil
}

func (b *metadataManager) GetBool(key EventMetadataKey) (bool, error) {
	iface, ok := b.Get(key)
	if !ok {
		return false, fmt.Errorf("not found")
	}

	boolValue, ok := iface.(bool)
	if !ok {
		return false, fmt.Errorf("invalid type")
	}
	return boolValue, nil
}

func (b *metadataManager) GetContext() context.Context {
	return context.WithValue(b.ctx, metadataKey, b.data)
}
