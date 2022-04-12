package m8.authz

alice_admin = {"User": {"Id": "1234", "Name": "alice", "Roles": [{"Name": "admin", "Scope": "system"}, {"Name": "admin", "Scope": "tenant", "Resource": "1234"}]}}

bob_tenant_admin = {"User": {"Id": "12345", "Name": "bob", "Roles": [{"Name": "admin", "Scope": "tenant"}]}}

jane = {"User": {"Id": "123456", "Name": "jane"}, "Path": "/domain.CommandHandlerExtensions/"}

scim_scope = {"Path": "/scim/something", "Authentication": {"Scopes": ["WRITE_SCIM"]}}

grpc_health_check = {"Path": "/grpc.health.v1.Health/Check"}

test_system_admin {
	is_system_admin with input as alice_admin
	not is_system_admin with input as bob_tenant_admin
}

test_authorized {
	authorized with input as alice_admin
	authorized with input as jane
	authorized with input as scim_scope
	authorized with input as grpc_health_check
}
