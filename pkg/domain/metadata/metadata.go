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

package metadata

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/internal/version"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/metadata"
)

const (
	componentName    = "component_name"
	componentVersion = "component_version"
	componentCommit  = "component_commit"
)

// UserInformation are identifying information about a user.
type UserInformation struct {
	Id        uuid.UUID
	Name      string
	Email     string
	NotBefore time.Time
}

// domainMetadataManager is a domain specific metadata manager.
type DomainMetadataManager struct {
	es.MetadataManager
	log logger.Logger
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (*DomainMetadataManager, error) {
	m := &DomainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
		logger.WithName("domain-metadata-manager"),
	}

	if len(m.GetMetadata()) == 0 {
		tags := grpc_ctxtags.Extract(ctx)

		// Get the grpc metadata from incoming context
		if tags != grpc_ctxtags.NoopTags {
			data := make(map[string]string)
			for k, v := range tags.Values() {
				m.log.V(logger.DebugLevel).Info("grpc metadata from incoming context", "key", k, "value", v)
				data[k] = v.(string)
			}
			m.SetMetadata(data)
		}
	}

	if _, exists := m.Get(componentName); !exists {
		m.SetComponentInformation()
	}

	return m, nil
}

// SetComponentInformation sets the ComponentInformation about the currently executing service/component.
func (m *DomainMetadataManager) SetComponentInformation() {
	m.Set(componentName, version.Name)
	m.Set(componentVersion, version.Version)
	m.Set(componentCommit, version.Commit)
}

// GetComponentName gets the name of the component which created the context.
func (m *DomainMetadataManager) GetComponentName() string {
	res, ok := m.Get(componentName)
	if ok {
		return res
	}
	return ""
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *DomainMetadataManager) SetUserInformation(userInformation *UserInformation) {
	m.Set(auth.HeaderAuthName, userInformation.Name)
	m.Set(auth.HeaderAuthEmail, userInformation.Email)
	m.Set(auth.HeaderAuthId, userInformation.Id.String())
	m.Set(auth.HeaderAuthNotBefore, userInformation.NotBefore.Format(auth.HeaderAuthNotBeforeFormat))
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *DomainMetadataManager) GetUserInformation() *UserInformation {
	userInfo := &UserInformation{}
	if header, ok := m.Get(auth.HeaderAuthName); ok {
		userInfo.Name = header
	}
	if header, ok := m.Get(auth.HeaderAuthEmail); ok {
		userInfo.Email = header
	}
	if header, ok := m.Get(auth.HeaderAuthId); ok {
		id, err := uuid.Parse(header)
		if err == nil {
			userInfo.Id = id
		}
	}
	if header, ok := m.Get(auth.HeaderAuthNotBefore); ok {
		t, err := time.Parse(auth.HeaderAuthNotBeforeFormat, header)
		if err == nil {
			userInfo.NotBefore = t
		}
	}
	return userInfo
}

// GetOutgoingGrpcContext returns a new context enriched with the metadata of this manager.
func (m *DomainMetadataManager) GetOutgoingGrpcContext() context.Context {
	return metadata.NewOutgoingContext(m.GetContext(), metadata.New(m.GetMetadata()))
}

func (m *DomainMetadataManager) GetContext() context.Context {
	return m.MetadataManager.GetContext()
}
