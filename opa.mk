
##@ OPA

opa-test: ## run all tests
	@opa test -v build/package/helm/gateway/files/policies.rego build/package/helm/gateway/files/policies_test.rego

opa-coverage: ## show coverage
	@opa test --coverage --format=json build/package/helm/gateway/files/policies.rego build/package/helm/gateway/files/policies_test.rego | jq .coverage
