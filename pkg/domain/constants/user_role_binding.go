package constants

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

// UserRoleBinding
const (
	UserRoleBindingType        es.AggregateType = "UserRoleBinding"
	CreateUserRoleBindingType  es.CommandType   = "CreateUserRoleBinding"
	UserRoleBindingCreatedType es.EventType     = "UserRoleBindingCreated"
)
