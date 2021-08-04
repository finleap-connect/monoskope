HELM                		?= helm3
HELM_OUTPUT_DIR             ?= tmp
HELM_REGISTRY 				?= https://artifactory.figo.systems/artifactory/virtual_helm
HELM_REGISTRY_ALIAS			?= finleap
HELM_RELEASE                ?= m8
KUBE_NAMESPACE ?= platform-monoskope-monoskope

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

helm-set-chart-version-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml version "$(VERSION)" --inplace

helm-set-app-version-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml appVersion "$(VERSION)" --inplace
	@yq write $(HELM_PATH)/$*/values.yaml image.tag "$(VERSION)" --inplace

helm-set-version-%:
	@$(MAKE) helm-set-chart-version-$*
	@$(MAKE) helm-set-app-version-$*

helm-set-app-version-latest-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml appVersion "$(LATEST_TAG)" --inplace
	@yq write $(HELM_PATH)/$*/values.yaml image.tag "$(LATEST_TAG)" --inplace

helm-set-version-all:
	@$(MAKE) helm-set-version-gateway
	@$(MAKE) helm-set-version-eventstore
	@$(MAKE) helm-set-version-commandhandler
	@$(MAKE) helm-set-version-queryhandler
	@$(MAKE) helm-set-version-cluster-bootstrap-reactor

helm-docs: ## update the auto generated docs of all helm charts
	@docker run --rm --volume "$(PWD):/helm-docs" -u $(shell id -u) jnorwood/helm-docs:v1.4.0 --template-files=./README.md.gotmpl
