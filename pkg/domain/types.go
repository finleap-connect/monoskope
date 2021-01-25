package domain

import (
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

const (
	UserRoleBinding AggregateType = "UserRoleBinding"
	AddRoleToUser   CommandType   = "AddRoleToUser"
	UserRoleAdded   EventType     = "UserRoleAdded"
)
