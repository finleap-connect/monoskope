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

.PHONY: helm-deploy helm-template helm-dependency-update

helm-template-clean:
	rm -R $(HELM_OUTPUT_DIR)

helm-dependency-update:
	$(HELM) dep update $(HELM_PATH_MONOSKOPE)

helm-deploy:
	$(HELM) upgrade --install monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

helm-template: helm-template-clean
	$(HELM) template monoskope $(HELM_PATH_MONOSKOPE) --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds