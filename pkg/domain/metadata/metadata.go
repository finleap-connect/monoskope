package metadata

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
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
		gateway.HeaderAuthEmail,
		gateway.HeaderAuthId,
		gateway.HeaderAuthIssuer,
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
	UserRoleBindings []*projections.UserRoleBinding
}

func newDomainContext(ctx *DomainContext) *DomainContext {
	if ctx == nil {
		return &DomainContext{}
	}

	return &DomainContext{
		UserRoleBindings: ctx.UserRoleBindings,
	}
}

// NewDomainMetadataManager creates a new domainMetadataManager to handle domain metadata via context.
func NewDomainMetadataManager(ctx context.Context) (*DomainMetadataManager, error) {
	m := &DomainMetadataManager{
		es.NewMetadataManagerFromContext(ctx),
		nil,
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
	m.Set(gateway.HeaderAuthName, userInformation.Name)
	m.Set(gateway.HeaderAuthEmail, userInformation.Email)
	m.Set(gateway.HeaderAuthIssuer, userInformation.Issuer)
	m.Set(gateway.HeaderAuthId, userInformation.Id.String())
}

// GetUserInformation returns the UserInformation stored in the metadata.
func (m *DomainMetadataManager) GetUserInformation() *UserInformation {
	userInfo := &UserInformation{}
	if header, ok := m.Get(gateway.HeaderAuthName); ok {
		userInfo.Name = header
	}
	if header, ok := m.Get(gateway.HeaderAuthEmail); ok {
		userInfo.Email = header
	}
	if header, ok := m.Get(gateway.HeaderAuthIssuer); ok {
		userInfo.Issuer = header
	}
	if header, ok := m.Get(gateway.HeaderAuthId); ok {
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
