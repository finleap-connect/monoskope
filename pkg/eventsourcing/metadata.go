package eventsourcing

import (
	"context"
	"encoding/json"
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
	GetMetadata() map[string][]byte
	// SetMetadata sets the metadata of this manager.
	SetMetadata(map[string][]byte) MetadataManager

	// Get returns the information stored for the key.
	Get(string) ([]byte, bool)
	GetObject(string, interface{}) error

	// Set stores information for the key.
	Set(string, []byte) MetadataManager
	SetObject(string, interface{}) error
}

type metadataManager struct {
	ctx  context.Context
	data map[string][]byte
}

func newMetadataManager() *metadataManager {
	return &metadataManager{
		ctx:  context.Background(),
		data: make(map[string][]byte),
	}
}

func NewMetadataManager() MetadataManager {
	return newMetadataManager()
}

func NewMetadataManagerFromContext(ctx context.Context) MetadataManager {
	m := newMetadataManager()
	m.ctx = ctx

	d, ok := ctx.Value(metadataKey).(map[string][]byte)
	if ok {
		m.data = d
	}

	return m
}

func (m *metadataManager) GetMetadata() map[string][]byte {
	return m.data
}

func (m *metadataManager) SetMetadata(metadata map[string][]byte) MetadataManager {
	m.data = metadata
	return m
}

func (m *metadataManager) Get(key string) ([]byte, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *metadataManager) Set(key string, value []byte) MetadataManager {
	m.data[key] = value
	return m
}

func (m *metadataManager) SetObject(key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.Set(key, bytes)
	return nil
}

func (m *metadataManager) GetObject(key string, v interface{}) error {
	bytes, found := m.Get(key)
	if !found {
		return fmt.Errorf("not found")
	}
	return json.Unmarshal(bytes, v)
}

// GetContext returns a new context enriched with the metadata of this manager.
func (m *metadataManager) GetContext() context.Context {
	return context.WithValue(m.ctx, metadataKey, m.data)
}
