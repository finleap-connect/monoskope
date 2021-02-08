package eventsourcing

import (
	"context"
	"fmt"
)

// EventMetadataKey is the base key type for metadata.
type EventMetadataKey int

// metadataKey is the key type for metadata stored in context.
var metadataKey EventMetadataKey

/*
 MetadataManager is an interface for a storage of metadata.
 It can be used to easily store any metadata in the context of a call.
*/
type MetadataManager interface {
	// GetContext returns a new context enriched with the metadata of this manager based on the .
	GetContext() context.Context

	// GetMetadata returns the metadata of this manager.
	GetMetadata() map[EventMetadataKey]interface{}
	// SetMetadata sets the metadata of this manager.
	SetMetadata(map[EventMetadataKey]interface{}) MetadataManager

	// Get returns the information stored for the key.
	Get(EventMetadataKey) (interface{}, bool)
	// Get returns the information stored for the key as bool.
	GetBool(EventMetadataKey) (bool, error)
	// Get returns the information stored for the key as string.
	GetString(EventMetadataKey) (string, error)

	// Set stores information for the key.
	Set(EventMetadataKey, interface{}) MetadataManager
}

type metadataManager struct {
	ctx  context.Context
	data map[EventMetadataKey]interface{}
}

func newMetadataManager() *metadataManager {
	return &metadataManager{
		ctx:  context.Background(),
		data: make(map[EventMetadataKey]interface{}),
	}
}

func NewMetadataManager() MetadataManager {
	return newMetadataManager()
}

func NewMetadataManagerFromContext(ctx context.Context) MetadataManager {
	m := newMetadataManager()
	m.ctx = ctx

	d, ok := ctx.Value(metadataKey).(map[EventMetadataKey]interface{})
	if ok {
		m.data = d
	}

	return m
}

func (m *metadataManager) GetMetadata() map[EventMetadataKey]interface{} {
	return m.data
}

func (m *metadataManager) SetMetadata(metadata map[EventMetadataKey]interface{}) MetadataManager {
	m.data = metadata
	return m
}

func (m *metadataManager) Get(key EventMetadataKey) (interface{}, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *metadataManager) Set(key EventMetadataKey, value interface{}) MetadataManager {
	m.data[key] = value
	return m
}

func (m *metadataManager) GetString(key EventMetadataKey) (string, error) {
	iface, ok := m.Get(key)
	if !ok {
		return "", fmt.Errorf("not found")
	}

	stringValue, ok := iface.(string)
	if !ok {
		return "", fmt.Errorf("invalid type")
	}
	return stringValue, nil
}

func (m *metadataManager) GetBool(key EventMetadataKey) (bool, error) {
	iface, ok := m.Get(key)
	if !ok {
		return false, fmt.Errorf("not found")
	}

	boolValue, ok := iface.(bool)
	if !ok {
		return false, fmt.Errorf("invalid type")
	}
	return boolValue, nil
}

// GetContext returns a new context enriched with the metadata of this manager.
func (m *metadataManager) GetContext() context.Context {
	return context.WithValue(m.ctx, metadataKey, m.data)
}
