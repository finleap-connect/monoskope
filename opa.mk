
##@ OPA

export POLICIES_PATH = $(BUILD_PATH)/build/package/helm/gateway/files/policies

opa-test: ## run all tests
	@opa test -v $(POLICIES_PATH)/policies.rego $(POLICIES_PATH)/policies_test.rego

opa-coverage: ## show coverage
	@opa test --coverage --format=json $(POLICIES_PATH)/policies.rego $(POLICIES_PATH)/policies_test.rego | jq .coverage
