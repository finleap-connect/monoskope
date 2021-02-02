package user

import (
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// User
const (
	UserType        AggregateType = "User"
	CreateUserType  CommandType   = "CreateUser"
	UserCreatedType EventType     = "UserCreated"
)

// UserRoleBinding
const (
	UserRoleBindingType        AggregateType = "UserRoleBinding"
	CreateUserRoleBindingType  CommandType   = "CreateUserRoleBinding"
	UserRoleBindingCreatedType EventType     = "UserRoleBindingCreated"
)
