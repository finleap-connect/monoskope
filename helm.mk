HELM                        ?= helm
HELM_OUTPUT_DIR             ?= tmp
HELM_REGISTRY               ?= https://finleap-connect.github.io/charts
HELM_REGISTRY_ALIAS         ?= finleap-connect
HELM_RELEASE                ?= m8
KUBE_NAMESPACE              ?= platform-monoskope
YQ							?= yq

.PHONY: template-clean dependency-update install uninstall template docs

##@ Helm

helm-clean: ## clean up templated helm charts
	@rm -Rf $(HELM_OUTPUT_DIR)

helm-dep-%: ## update helm dependencies
	@$(HELM) dep update $(HELM_PATH)/$*

helm-lint-%: ## lint helm chart
	@$(HELM) lint $(HELM_PATH)/$*

helm-install-%: ## install helm chart from local sources
	@cat $(HELM_VALUES_FILE) | sed "s/0.0.1-local/$(VERSION)/g" > $(HELM_VALUES_FILE).tag
	@$(HELM) upgrade --install $(HELM_RELEASE) $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE).tag --skip-crds
	@rm $(HELM_VALUES_FILE).tag

helm-install-from-repo-%: ## install helm chart from build artifact
	@$(HELM) repo update
	@$(HELM) upgrade --install $(HELM_RELEASE) $(HELM_REGISTRY_ALIAS)/$* --namespace $(KUBE_NAMESPACE) --version $(VERSION) --values $(HELM_VALUES_FILE) --skip-crds

helm-uninstall-%: ## uninstall helm chart
	@$(HELM) uninstall $(HELM_RELEASE) --namespace $(KUBE_NAMESPACE)

helm-template-%: helm-clean ## template helm chart
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template $(HELM_RELEASE) $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds
	@echo "ATTENTION:"
	@echo "If you want to have the latest dependencies (e.g. gateway chart changes)"
	@echo "execute the following command prior to the current command:"
	@echo "$$ $(MAKE) helm-dep-$*"
	@echo

helm-add-finleap: ## add finleap helm chart repo
	@$(HELM) repo add --username $(HELM_USER) --password $(HELM_PASSWORD) $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"

helm-set-version-all:
	@find $(HELM_PATH) -name 'Chart.yaml' -exec $(YQ) e --inplace '.version = "$(VERSION)"' {} \;
	@find $(HELM_PATH) -name 'Chart.yaml' -exec $(YQ) e --inplace '.appVersion = "$(VERSION)"' {} \;
	@find $(HELM_PATH) -name 'Chart.yaml' -exec $(YQ) e --inplace '(.dependencies.[].version | select(. == "0.0.1-local")) |= "$(VERSION)"' {} \;

helm-docs: ## update the auto generated docs of all helm charts
	@docker run --rm --volume "$(PWD):/helm-docs" -u $(shell id -u) jnorwood/helm-docs:v1.4.0 --template-files=./README.md.gotmpl
