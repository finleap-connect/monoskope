package m8.authz

import future.keywords.in

default authorized = false

allowed_paths = ["/domain.CommandHandlerExtensions/", "/domain.Cluster/GetByName"]

authorized {
	some path in allowed_paths
	startswith(input.Path, path)
}

authorized {
	is_system_admin
}

is_system_admin {
	some role in input.User.Roles
	"system" == role.Scope
	"admin" == role.Name
}
