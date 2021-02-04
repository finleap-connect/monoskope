package aggregates

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

const (
	// Type for the UserAggregate
	User es.AggregateType = "User"
	// Type for the UserRoleBindingAggregate
	UserRoleBinding es.AggregateType = "UserRoleBinding"
)
