package aggregates

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Type for the UserAggregate
	User es.AggregateType = "User"
	// Type for the UserRoleBindingAggregate
	UserRoleBinding es.AggregateType = "UserRoleBinding"
)
