package m8.authz

import future.keywords.in

default authorized = false

# list of allowed paths for any user
allowed_paths := [
	"/grpc.health.v1.Health/",
	"/gateway.Gateway/",
	"/gateway.ClusterAuth/GetAuthToken/",
	"/domain.CommandHandlerExtensions/",
	"/domain.Cluster/GetByName",
]

scoped_paths := [{"path": "/scim/", "scope": "WRITE_SCIM"}]

# check if system admin
is_system_admin {
	some role in input.User.Roles
	role.Scope == "system"
	role.Name == "admin"
}

# authorized because system admin
authorized {
	is_system_admin
}

# authorized via allowed_paths
authorized {
	some path in allowed_paths
	startswith(input.Path, path)
}

# authorized via scope
authorized {
	some scoped_path in scoped_paths
	startswith(input.Path, scoped_path.path)
	some scope in input.Authentication.Scopes
	scope == scoped_path.scope
}
