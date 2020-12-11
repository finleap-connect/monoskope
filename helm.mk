HELM                		?= helm3
HELM_OUTPUT_DIR             ?= tmp
HELM_REGISTRY 				?= https://artifactory.figo.systems/artifactory/virtual_helm
HELM_REGISTRY_ALIAS			?= finleap


.PHONY: helm-template-clean helm-dependency-update helm-install helm-uninstall helm-template

clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

dep-%:
	@$(HELM) dep update $(HELM_PATH)/$*

lint-%:
	@$(HELM) lint $(HELM_PATH)/$*

install-%:
	@$(MAKE) helm-dep-$*
	@$(HELM) upgrade --install m8dev $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

install-from-repo-%:
	@$(HELM) repo update
	@$(HELM) upgrade --install m8dev $(HELM_REGISTRY_ALIAS)/$* --namespace $(KUBE_NAMESPACE) --version $(VERSION) --values $(HELM_VALUES_FILE)

uninstall-%: 
	@$(HELM) uninstall m8dev --namespace $(KUBE_NAMESPACE)

template-%: clean
	@$(MAKE) helm-dep-$*
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template m8dev $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug

add-kubism:
	@$(HELM) repo add kubism.io https://kubism.github.io/charts/

add-finleap:
	@$(HELM) repo add --username $(HELM_USER) --password $(HELM_PASSWORD) $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"

set-chart-version-%:
	yq write $(HELM_PATH)/$*/Chart.yaml version "$(VERSION)" --inplace

set-app-version-%:
	yq write $(HELM_PATH)/$*/Chart.yaml appVersion "$(VERSION)" --inplace
	yq write $(HELM_PATH)/$*/values.yaml image.tag "$(VERSION)" --inplace

set-version-%:
	@$(MAKE) helm-set-chart-version-$*
	@$(MAKE) helm-set-app-version-$*
