package aggregates

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Type for the UserAggregate
	User es.AggregateType = "User"
	// Type for the UserRoleBindingAggregate
	UserRoleBinding es.AggregateType = "UserRoleBinding"
	// Type for the TenantAggregate
	Tenant es.AggregateType = "Tenant"
	// Type for the ClusterRegistrationAggregate
	ClusterRegistration es.AggregateType = "ClusterRegistration"
	// Type for the ClusterAggregate
	Cluster es.AggregateType = "Cluster"
)
