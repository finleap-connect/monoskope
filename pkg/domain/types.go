package domain

import (
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

// User
const (
	User        AggregateType = "User"
	CreateUser  CommandType   = "CreateUser"
	UserCreated EventType     = "UserCreated"
)

// UserRoleBinding
const (
	UserRoleBinding        AggregateType = "UserRoleBinding"
	CreateUserRoleBinding  CommandType   = "CreateUserRoleBinding"
	UserRoleBindingCreated EventType     = "UserRoleBindingCreated"
)
