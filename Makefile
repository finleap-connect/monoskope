VERSION   ?= 0.0.1-local


HELM                ?= helm3
HELM_PATH_MONOSKOPE ?= build/package/helm/monoskope
HELM_VALUES_FILE    ?= examples/00-monoskope-values.yaml
HELM_OUTPUT_DIR     ?= tmp

KUBE_NAMESPACE ?= platform-monoskope-monoskope

# go

go-%:
	$(MAKE) -f go.mk $*

# helm

.PHONY: helm-template-clean helm-dependency-update helm-install helm-uninstall helm-template

helm-template-clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

helm-dependency-update:
	@$(HELM) dep update $(HELM_PATH_MONOSKOPE)

helm-install:
	@$(HELM) upgrade --install monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

helm-uninstall:
	@$(HELM) uninstall monoskope --namespace $(KUBE_NAMESPACE)

helm-template: helm-template-clean
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds