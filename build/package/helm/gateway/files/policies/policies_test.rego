package m8.authz

alice_admin = {"user": {"name": "alice", "roles": [{"name": "admin", "scope": "system"}, {"name": "admin", "scope": "tenant", "resource": "1234"}]}}

bob_tenant_admin = {"user": {"name": "bob", "roles": [{"name": "admin", "scope": "tenant"}]}}

test_system_admin {
	is_system_admin with input as alice_admin
	not is_system_admin with input as bob_tenant_admin
}

test_authorized {
	authorized with input as alice_admin
}
