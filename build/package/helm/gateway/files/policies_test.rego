package m8.authz

test_allow_with_data {
	allow with input as {"user": {"name": "alice", "roles": [{"name": "admin", "scope": "system"}]}}
}
