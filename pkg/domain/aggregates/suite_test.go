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

package aggregates

import (
	"context"
	"testing"

	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/commands"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	meta "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	expectedUserName = "the one cluster"
	expectedEmail    = "me@example.com"

	expectedTenantScope = common.Scope_tenant
	expectedAdminRole   = common.Role_admin
	expectedResourceId  = uuid.New()
	expectedUserId      = uuid.New()

	expectedCSR                     = []byte("This should be a CSR")
	expectedReferencedAggregateId   = uuid.New()
	expectedReferencedAggregateType = aggregates.Cluster
	expectedCertificate             = []byte("this should be the certificate")

	expectedClusterDisplayName      = "the one cluster"
	expectedClusterName             = "one-cluster"
	expectedClusterApiServerAddress = "one.example.com"
	expectedClusterCACertBundle     = []byte("This should be a certificate")
)

func TestAggregates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aggregates Suite")
}

func createSysAdminCtx() context.Context {
	metaMgr, err := meta.NewDomainMetadataManager(context.Background())
	Expect(err).NotTo(HaveOccurred())

	metaMgr.SetUserInformation(&meta.UserInformation{
		Id:     uuid.New(),
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})

	metaMgr.SetRoleBindings([]*projections.UserRoleBinding{
		{
			Role:  common.Role_admin.String(),
			Scope: common.Scope_system.String(),
		},
	})

	return metaMgr.GetContext()
}

func createCluster(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateClusterCommand(uuid.New()).(*cmd.CreateClusterCommand)
	Expect(ok).To(BeTrue())

	esCommand.CreateCluster.DisplayName = expectedClusterDisplayName
	esCommand.CreateCluster.Name = expectedClusterName
	esCommand.CreateCluster.ApiServerAddress = expectedClusterApiServerAddress
	esCommand.CreateCluster.CaCertBundle = expectedClusterCACertBundle

	return agg.HandleCommand(ctx, esCommand)
}

func newRequestCertificateCommand(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewRequestCertificateCommand(uuid.New()).(*commands.RequestCertificateCommand)
	Expect(ok).To(BeTrue())

	esCommand.SigningRequest = expectedCSR
	esCommand.ReferencedAggregateId = expectedReferencedAggregateId.String()
	esCommand.ReferencedAggregateType = expectedReferencedAggregateType.String()

	return agg.HandleCommand(ctx, esCommand)
}

type aggregateTestStore struct {
	bindings map[uuid.UUID]es.Aggregate
	users    map[uuid.UUID]es.Aggregate
}

// NewTestAggregateManager creates a new dummy AggregateHandler which allows observing interactions and injecting test data.
func NewTestAggregateManager() es.AggregateStore {
	return &aggregateTestStore{
		bindings: make(map[uuid.UUID]es.Aggregate),
		users:    make(map[uuid.UUID]es.Aggregate),
	}
}

func (tas *aggregateTestStore) Add(agg es.Aggregate) {
	switch agg.Type() {
	case aggregates.User:
		tas.users[agg.ID()] = agg
	case aggregates.UserRoleBinding:
		tas.bindings[agg.ID()] = agg
	}
}

// Get returns the most recent version of all aggregate of a given type.
func (tas *aggregateTestStore) All(ctx context.Context, atype es.AggregateType) ([]es.Aggregate, error) {
	var retmap map[uuid.UUID]es.Aggregate

	switch atype {
	case aggregates.Certificate:
		return []es.Aggregate{}, nil
	case aggregates.Cluster:
		return []es.Aggregate{}, nil
	case aggregates.Tenant:
		return []es.Aggregate{}, nil
	case aggregates.UserRoleBinding:
		retmap = tas.bindings
	case aggregates.User:
		retmap = tas.users
	}

	values := make([]es.Aggregate, 0, len(retmap))
	for _, aggr := range retmap {
		values = append(values, aggr)
	}

	return values, nil
}

// Get returns the most recent version of an aggregate.
func (tas *aggregateTestStore) Get(ctx context.Context, atype es.AggregateType, id uuid.UUID) (es.Aggregate, error) {
	var (
		retmap      map[uuid.UUID]es.Aggregate
		notFoundVal error
	)

	switch atype {
	case aggregates.Certificate:
		return nil, nil // this will break if the aggregate is actually used. Implement if needed
	case aggregates.Cluster:
		return nil, nil // this will break if the aggregate is actually used. Implement if needed
	case aggregates.Tenant:
		return nil, nil // this will break if the aggregate is actually used. Implement if needed
	case aggregates.UserRoleBinding:
		retmap = tas.bindings
		notFoundVal = errors.ErrUserRoleBindingNotFound
	case aggregates.User:
		retmap = tas.users
		notFoundVal = errors.ErrUserNotFound
	default:
		return nil, errors.ErrUnknownAggregateType
	}

	ret := retmap[id]
	if ret == nil {
		return nil, notFoundVal
	}
	return ret, nil

}

// Update stores all in-flight events for an aggregate.
func (tas *aggregateTestStore) Update(context.Context, es.Aggregate) error {
	return nil
}

func createTenant(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateTenantCommand(agg.ID()).(*cmd.CreateTenantCommand)
	Expect(ok).To(BeTrue())

	esCommand.CreateTenantCommandData.Name = expectedTenantName
	esCommand.CreateTenantCommandData.Prefix = expectedPrefix

	return agg.HandleCommand(ctx, esCommand)
}

func createUserRoleBinding(ctx context.Context, agg es.Aggregate, userId uuid.UUID) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateUserRoleBindingCommand(uuid.New()).(*cmd.CreateUserRoleBindingCommand)
	Expect(ok).To(BeTrue())

	esCommand.UserId = userId.String()
	esCommand.Role = expectedAdminRole
	esCommand.Scope = expectedTenantScope
	esCommand.Resource = expectedResourceId.String()

	return agg.HandleCommand(ctx, esCommand)
}

func createUser(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateUserCommand(uuid.New()).(*cmd.CreateUserCommand)
	Expect(ok).To(BeTrue())

	esCommand.Name = expectedUserName
	esCommand.Email = expectedEmail

	return agg.HandleCommand(ctx, esCommand)
}
