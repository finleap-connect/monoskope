package eventsourcing

import (
	"context"
	"fmt"
)

// metadataKey is the key type for metadata stored in context.
type metadataKeyType struct {
}

var metadataKey metadataKeyType

/*
 MetadataManager is an interface for a storage of metadata.
 It can be used to easily store any metadata in the context of a call.
*/
type MetadataManager interface {
	// GetContext returns a new context enriched with the metadata of this manager based on the .
	GetContext() context.Context

	// GetMetadata returns the metadata of this manager.
	GetMetadata() map[string]interface{}
	// SetMetadata sets the metadata of this manager.
	SetMetadata(map[string]interface{}) MetadataManager

	// Get returns the information stored for the key.
	Get(string) (interface{}, bool)
	// Get returns the information stored for the key as bool.
	GetBool(string) (bool, error)
	// Get returns the information stored for the key as string.
	GetString(string) (string, error)

	// Set stores information for the key.
	Set(string, interface{}) MetadataManager
}

type metadataManager struct {
	ctx  context.Context
	data map[string]interface{}
}

func newMetadataManager() *metadataManager {
	return &metadataManager{
		ctx:  context.Background(),
		data: make(map[string]interface{}),
	}
}

func NewMetadataManager() MetadataManager {
	return newMetadataManager()
}

func NewMetadataManagerFromContext(ctx context.Context) MetadataManager {
	m := newMetadataManager()
	m.ctx = ctx

	d, ok := ctx.Value(metadataKey).(map[string]interface{})
	if ok {
		m.data = d
	}

	return m
}

func (m *metadataManager) GetMetadata() map[string]interface{} {
	return m.data
}

func (m *metadataManager) SetMetadata(metadata map[string]interface{}) MetadataManager {
	m.data = metadata
	return m
}

func (m *metadataManager) Get(key string) (interface{}, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *metadataManager) Set(key string, value interface{}) MetadataManager {
	m.data[key] = value
	return m
}

func (m *metadataManager) GetString(key string) (string, error) {
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

func (m *metadataManager) GetBool(key string) (bool, error) {
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
