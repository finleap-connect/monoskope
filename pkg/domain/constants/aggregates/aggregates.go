// Copyright 2022 Monoskope Authors
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

import es "github.com/finleap-connect/monoskope/pkg/eventsourcing"

const (
	// Type for the UserAggregate
	User es.AggregateType = "User"
	// Type for the UserRoleBindingAggregate
	UserRoleBinding es.AggregateType = "UserRoleBinding"
	// Type for the TenantAggregate
	Tenant es.AggregateType = "Tenant"
	// Type for the ClusterAggregate
	Cluster es.AggregateType = "Cluster"
	// Type for the CertificateAggregate
	Certificate es.AggregateType = "Certificate"
	// Type for the TenantClusterBindingAggregate
	TenantClusterBinding es.AggregateType = "TenantClusterBinding"
)
