package m8.authz

import future.keywords.in

default allow = false

allow {
	is_admin
}

is_admin {
	some role in input.user.roles
	"system" == role.scope
	"admin" == role.name
}
