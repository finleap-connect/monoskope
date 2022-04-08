package m8.authz

import future.keywords.in

default authorized = false

authorized {
	startswith(input.Path, "/domain.CommandHandlerExtensions/")
}

authorized {
	is_system_admin
}

is_system_admin {
	some role in input.User.Roles
	"system" == role.Scope
	"admin" == role.Name
}
