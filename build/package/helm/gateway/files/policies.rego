package m8.authz

import future.keywords.in

default allow = false

allow {
	is_system_admin
}

is_system_admin {
	some role in input.user.roles
	"system" == role.scope
	"admin" == role.name
}
