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

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/version"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/grpc/metadata"
)

const (
	componentName    = "component_name"
	componentVersion = "component_version"
	componentCommit  = "component_commit"
)

var (
	acceptedHeaders = []string{
		componentName,
		componentCommit,
		componentVersion,
		auth.HeaderAuthEmail,
		auth.HeaderAuthId,
		auth.HeaderAuthIssuer,
	}
)

// UserInformation are identifying information about a user.
type UserInformation struct {
	Id     uuid.UUID
	Name   string
	Email  string
	Issuer string
}

// domainMetadataManager is a domain specific metadata manager.
type DomainMetadataManager struct {
	es.MetadataManager
	domainContext *DomainContext
}

type DomainContext struct {
	context.Context
	UserRoleBindings    []*projections.UserRoleBinding
	BypassAuthorization bool
}

func newDomainContext(ctx *DomainContext) *DomainContext {
	if ctx == nil {
		return &DomainContext{}
	}

	return &DomainContext{
		UserRoleBindings:    ctx.UserRoleBindings,
		BypassAuthorization: ctx.BypassAuthorization,
	}
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (*DomainMetadataManager, error) {
	m := &DomainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
		newDomainContext(nil),
	}

	if len(m.GetMetadata()) == 0 {
		// Get the grpc metadata from incoming context
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			data := make(map[string]string)
			for k, v := range md {
				if isHeaderAccepted(k) {
					data[k] = v[0] // typically only the first and only value of that is relevant
				}
			}
			m.SetMetadata(data)
		}
	}

	if domainContext, ok := ctx.(*DomainContext); ok {
		m.domainContext = domainContext
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

func (m *DomainMetadataManager) SetRoleBindings(roleBindings []*projections.UserRoleBinding) {
	m.domainContext.UserRoleBindings = roleBindings
}

func (m *DomainMetadataManager) GetRoleBindings() []*projections.UserRoleBinding {
	return m.domainContext.UserRoleBindings
}

// SetUserInformation sets the UserInformation in the metadata.
func (m *DomainMetadataManager) SetUserInformation(userInformation *UserInformation) {
	m.Set(auth.HeaderAuthName, userInformation.Name)
	m.Set(auth.HeaderAuthEmail, userInformation.Email)
	m.Set(auth.HeaderAuthIssuer, userInformation.Issuer)
	m.Set(auth.HeaderAuthId, userInformation.Id.String())
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
	if header, ok := m.Get(auth.HeaderAuthIssuer); ok {
		userInfo.Issuer = header
	}
	if header, ok := m.Get(auth.HeaderAuthId); ok {
		id, err := uuid.Parse(header)
		if err == nil {
			userInfo.Id = id
		}
	}
	return userInfo
}

// GetOutgoingGrpcContext returns a new context enriched with the metadata of this manager.
func (m *DomainMetadataManager) GetOutgoingGrpcContext() context.Context {
	return metadata.NewOutgoingContext(m.GetContext(), metadata.New(m.GetMetadata()))
}

func (m *DomainMetadataManager) GetContext() context.Context {
	dc := newDomainContext(m.domainContext)
	dc.Context = m.MetadataManager.GetContext()
	return dc
}

func isHeaderAccepted(key string) bool {
	for _, acceptedHeader := range acceptedHeaders {
		if acceptedHeader == key {
			return true
		}
	}
	return false
}

// BypassAuthorization disables authorization checks and returns a function to enable it again
func (m *DomainMetadataManager) BypassAuthorization() func() {
	dc := newDomainContext(m.domainContext)
	dc.BypassAuthorization = true
	m.domainContext = dc
	return func() {
		dc.BypassAuthorization = false
	}
}

func (m *DomainMetadataManager) IsAuthorizationBypassed() bool {
	return m.domainContext.BypassAuthorization
}
