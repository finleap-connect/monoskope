package m8.authz

test_allow_system_admin {
	is_system_admin with input as {"user": {"name": "alice", "roles": [{"name": "admin", "scope": "system"}]}}
}

test_allow_system_admin {
	not is_system_admin with input as {"user": {"name": "alice", "roles": [{"name": "admin", "scope": "tenant"}]}}
}
