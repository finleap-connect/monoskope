package m8.authz

import future.keywords.in

default authorized = false

# list of allowed paths for any user
allowed_paths := [
	"/gateway.Gateway/",
	"/gateway.ClusterAuth/GetAuthToken/",
	"/domain.CommandHandlerExtensions/",
	"/domain.Cluster/GetByName",
]

scoped_paths := [{"path": "/scim/", "scope": "WRITE_SCIM"}]

command_path := "/eventsourcing.CommandHandler/Execute"

scope_system = "system"

scope_tenant = "tenant"

role_admin = "admin"

# check if system admin
is_system_admin {
	some role in input.User.Roles
	role.Scope == scope_system
	role.Name == role_admin
	print(input.User.Name, "is system admin")
}

# check if user is tenant admin and adjusts rolebindings of other users of the tenant
tenant_admin_rolebindings {
	# check that it is a command
	startswith(input.Path, command_path)
	req := json.unmarshal(input.Request)

	# check that it is related to user role bindings
	some type in input.CommandTypes.UserRoleBindingTypes
	req.type == type

	# check that target scope is tenant
	req.data.scope == scope_tenant

	# check that user is admin for the same tenant 
	some role in input.User.Roles
	role.Scope == scope_tenant
	role.Name == role_admin
	role.Resource == req.data.resource

	print(input.User.Name, "is tenant admin and allowed to execute", req.type)
}

# authorized because system admin
authorized {
	is_system_admin
}

# authorized because explicit rule
authorized {
	tenant_admin_rolebindings
}

# authorized via allowed_paths
authorized {
	some path in allowed_paths
	startswith(input.Path, path)
	print(path, "is allowed to everyone")
}

# authorized via scope
authorized {
	some scoped_path in scoped_paths
	startswith(input.Path, scoped_path.path)
	some scope in input.Authentication.Scopes
	scope == scoped_path.scope
	print("scope", scope, "allows access to path", scoped_path.path)
}
