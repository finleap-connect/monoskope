HELM                ?= helm3
HELM_PATH_MONOSKOPE ?= build/package/helm/monoskope
HELM_VALUES_FILE    ?= examples/00-monoskope-dev-values.yaml
HELM_OUTPUT_DIR     ?= tmp

.PHONY: helm-template-clean helm-dependency-update helm-install helm-uninstall helm-template

template-clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

dep:
	@$(HELM) dep update $(HELM_PATH_MONOSKOPE)

lint:
	@$(HELM) lint $(HELM_PATH_MONOSKOPE)

install: lint
	@$(HELM) upgrade --install monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

uninstall: 
	@$(HELM) uninstall monoskope --namespace $(KUBE_NAMESPACE)

template: helm-template-clean lint
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug