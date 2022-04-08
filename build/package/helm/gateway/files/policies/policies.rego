package m8.authz

import future.keywords.in

default authorized = false

# list of allowed paths for any user
allowed_paths = [
	"/gateway.Gateway/",
	"/gateway.ClusterAuth/GetAuthToken/",
	"/domain.CommandHandlerExtensions/",
	"/domain.Cluster/GetByName",
]

# check if system admin
is_system_admin {
	some role in input.User.Roles
	"system" == role.Scope
	"admin" == role.Name
}

# authorized via allowed_paths
authorized {
	some path in allowed_paths
	startswith(input.Path, path)
}

# authorized because system admin
authorized {
	is_system_admin
}
