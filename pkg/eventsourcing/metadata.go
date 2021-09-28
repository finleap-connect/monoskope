// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventsourcing

import (
	"context"
	"encoding/json"
	"fmt"
)

// metadataKey is the key type for metadata stored in context.
type metadataKeyType struct {
}

/*
 MetadataManager is an interface for a storage of metadata.
 It can be used to easily store any metadata in the context of a call.
*/
type MetadataManager interface {
	// GetContext returns a new context enriched with the metadata of this manager.
	GetContext() context.Context

	// GetMetadata returns the metadata of this manager.
	GetMetadata() map[string]string

	// SetMetadata sets the metadata of this manager.
	SetMetadata(map[string]string) MetadataManager

	// Get returns the information stored for the key.
	Get(string) (string, bool)
	GetObject(string, interface{}) error

	// Set stores information for the key.
	Set(string, string) MetadataManager
	SetObject(string, interface{}) error
}

type metadataManager struct {
	ctx  context.Context
	data map[string]string
}

func NewMetadataManagerFromContext(ctx context.Context) MetadataManager {
	m := &metadataManager{
		ctx:  ctx,
		data: make(map[string]string),
	}

	d, ok := ctx.Value(metadataKeyType{}).(map[string]string)
	if ok {
		m.data = d
	}

	return m
}

func (m *metadataManager) GetMetadata() map[string]string {
	return m.data
}

func (m *metadataManager) SetMetadata(metadata map[string]string) MetadataManager {
	m.data = metadata
	return m
}

func (m *metadataManager) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *metadataManager) Set(key string, value string) MetadataManager {
	m.data[key] = value
	return m
}

func (m *metadataManager) SetObject(key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.Set(key, string(bytes))
	return nil
}

func (m *metadataManager) GetObject(key string, v interface{}) error {
	jsonString, found := m.Get(key)
	if !found {
		return fmt.Errorf("metadata for key %s not found", key)
	}
	return json.Unmarshal([]byte(jsonString), v)
}

// GetContext returns a new context enriched with the metadata of this manager.
func (m *metadataManager) GetContext() context.Context {
	return context.WithValue(m.ctx, metadataKeyType{}, m.data)
}
