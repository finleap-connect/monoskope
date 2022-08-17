package m8.authz

import future.keywords.in

default authorized = false

# list of allowed paths for any user
allowed_paths := [
	"/gateway.Gateway",
	"/gateway.ClusterAuth/GetAuthToken",
	"/domain.CommandHandlerExtensions",
	"/common.ServiceInformationService/GetServiceInformation",
	"/domain.ClusterAccess/GetClusterAccess",
	"/domain.User",
	"/domain.Tenant",
	"/domain.Cluster",
]

scoped_paths := [{"scope": "WRITE_SCIM", "paths": [
	"/scim/",
	"/eventsourcing.CommandHandler/Execute",
	"/domain.User/",
]}]

command_path := "/eventsourcing.CommandHandler/Execute"

scope_system = "system"

scope_tenant = "tenant"

role_admin = "admin"

# check if system admin
is_system_admin {
	print("entering is_system_admin")
	some role in input.User.Roles
	role.Scope == scope_system
	role.Name == role_admin
	print(input.User.Name, "is system admin")
}

# check if user is tenant admin and adjusts rolebindings of other users of the tenant
tenant_admin_rolebindings {
	print("entering tenant_admin_rolebindings")

	# check that it is a command
	startswith(input.Path, command_path)
	req := json.unmarshal(input.Request)

	# check that it is related to user role bindings
	some type in input.CommandTypes.UserRoleBinding
	req.type == type

	# check that target scope is tenant
	req.data.scope == scope_tenant

	# check that user is admin for the same tenant 
	some role in input.User.Roles
	role.Scope == scope_tenant
	role.Name == role_admin
	role.Resource == req.data.resource

	print(input.User.Name, "is tenant admin and allowed to execute", req.type, "for tenant", req.data.resource)
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
	print("entering allowed_paths")
	some path in allowed_paths
	startswith(input.Path, path)
	print(path, "is allowed to everyone")
}

# authorized via scope
authorized {
	print("entering scoped_paths")
	some scoped_path in scoped_paths
	some scope in input.Authentication.Scopes
	scope == scoped_path.scope
	some path in scoped_path.paths
	startswith(input.Path, path)
	print("scope", scope, "allows access to path", path)
}
